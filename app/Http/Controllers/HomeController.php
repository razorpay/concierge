<?php

namespace App\Http\Controllers;

use App;
use Auth;
use Mail;
use View;
use DB;
use SocialOAuth;
use Request;
use Response;
use App\Models;
use Validator;
use Exception;
use Carbon\Carbon;
use Maclof\Kubernetes\Client;
use Illuminate\Support\Facades\Log;

class HomeController extends BaseController
{
    // 6 hours
    const MAX_EXPIRY = 32400;

    const CONCIERGE_TAG = 'concierge';

    /**
     * Stage One - The Login form.
     *
     * @return Login View
     */
    public function getIndex()
    {
        $code = Request::get('code');

        $google_service = SocialOAuth::consumer('Google', config('app.url'));

        if (! $code) {
            $url = $google_service->getAuthorizationUri([
                'hosted_domain' =>  config('concierge.google_domain')
            ]);

            return redirect((string) $url);
        } else {
            $token = $google_service->requestAccessToken($code);

            $response = $google_service->request(config('oauth-5-laravel.userinfo_url'));
            $result = json_decode($response);

            // Email must be:
            // - Verified
            // - Belong to razorpay.com domain
            //
            // Then only we'll create a user entry in the system or check for one

            if ($result->verified_email !== true or $result->hd !== config('concierge.google_domain'))
            {
                return App::abort(404);
            }

            // Find the user by email
            $user = Models\User::where('email', $result->email)->first();

            if ($user) {
                // Update some fields
                $user->access_token = $token->getAccessToken();
                $user->google_id = $result->id;
                $user->password = ''; // backward compatibility

                $user->save();

                // Login the user into the app
                Auth::loginUsingId($user->id);

                $redirectUrl = config('app.url') . '/groups';
                return redirect($redirectUrl);

            } else {
                App::abort(401);
            }
        }
    }

    /**
     * Log user out.
     *
     * @return Redirect to Home
     */
    public function getLogout()
    {
        Auth::logout();

        return redirect('/');
    }

    /**
     * Get the list of all security groups.
     *
     * @return getGroups view
     */
    public function getGroups()
    {
        //Get All security groups
        $ec2 = App::make('aws')->createClient('ec2');

        $filters = [
            [
                "Name"   => "tag-key",
                "Values" => [self::CONCIERGE_TAG]
            ],[
                "Name"   => "tag-value",
                "Values" => ["true"]
            ],
        ];

        $security_groups = $ec2->describeSecurityGroups(
            ['Filters' => $filters]
        );

        $security_groups = $security_groups['SecurityGroups'];

        error_reporting(E_ALL ^ E_NOTICE ^ E_WARNING);

        //Get all active leases
        $leases = Models\Lease::where('lease_type', 'aws')->get();

        //get all active Invites
        $invites = Models\Invite::get();

        return view('getGroups', [
            'security_groups'   => $security_groups,
            'leases'            => $leases,
            'invites'           => $invites,
        ]);
    }

    public function getIngresses()
    {
        $client = new Client([
            'master'        => 'http://127.0.0.1:8001',
        ]);

        $ingresses = $client->sendRequest('GET', '/ingresses', [
            'labelSelector' => 'concierge=true'
        ], null, false,  'extensions/v1beta1')['items'];


        array_walk($ingresses, function(&$ingress) {
            if (!isset($ingress['spec']['rules'])) {
                $ingress['hosts'] = [];
            }
            else {
                $hosts = [];
                foreach ($ingress['spec']['rules'] as $rule) {
                    if (isset($rule['host'])) {
                        $hosts[] = $rule['host'];
                    }
                }
                $ingress['hosts'] = $hosts;
            }
        });

        return view('kubernetes', [
            'ingresses'   => $ingresses,
        ]);
    }

    /*
     * Displays a security groups details with active leases & security rules.
     * @return getManage View
     */
    public function getManage($group_id)
    {

        //get security group details
        $ec2 = App::make('aws')->createClient('ec2');
        $security_group = $ec2->describeSecurityGroups([
            'GroupIds' => [$group_id],
        ]);
        $security_group = $security_group['SecurityGroups'][0];

        //get Active Leases
        $leases = Models\Lease::getByGroupId($group_id);

        //get Active Invites
        $invites = Models\Invite::getByGroupId($group_id);

        return View::make('getManage')
                    ->with('security_group', $security_group)
                    ->with('leases', $leases)
                    ->with('invites', $invites);
    }

    /*
     * Handles Lease creation & termination post requests to Group Manage page
     * @return Redirect to getManage View with error/success
     */
    public function postManage($group_id)
    {
        $input = Request::all();
        $messages = [];
        $email = null;
        /*
         For Lease Creation
        */
        if ('ssh' == $input['rule_type']) {
            $protocol = 'tcp';
            $port_from = '22';
            $port_to = '22';
        } elseif ('https' == $input['rule_type']) {
            $protocol = 'tcp';
            $port_from = '443';
            $port_to = '443';
        } elseif ('custom' == $input['rule_type']) {
            $protocol = $input['protocol'];
            $port_from = $input['port_from'];
            $port_to = $input['port_to'];

            //Validations
            if ($protocol != 'tcp' && $protocol != 'udp') {
                array_push($messages, 'Invalid Protocol');
            }
            if (! is_numeric($port_from) || $port_from > 65535 || $port_from <= 0) {
                array_push($messages, 'Invalid From port');
            }
            if (! is_numeric($port_to) || $port_to > 65535 || $port_to <= 0) {
                array_push($messages, 'Invalid To port');
            }
            if ($port_from > $port_to) {
                array_push($messages, 'From port Must be less than equal to To Port');
            }
        } else {
            App::abort(403, 'Unauthorized action.');
        }

        //Other validations
        $expiry = $input['expiry'];
        if (! is_numeric($expiry) or $expiry <= 0 or $expiry > self::MAX_EXPIRY) {
            array_push($messages, 'Invalid Expiry Time');
        }
        if (! in_array($input['access'], [1, 2, 3, 4])) {
            array_push($messages, 'Invalid invite Email');
        }
        if (2 == $input['access']) {
            if (! isset($input['email']) || ! filter_var($input['email'], FILTER_VALIDATE_EMAIL)) {
                array_push($messages, 'Invalid invite Email');
            }
        }

        //Validation fails
        if (! empty($messages)) {
            return redirect("/manage/$group_id")
                            ->with('message', implode('<br/>', $messages));
        }

        if (1 == $input['access']) {
            //Creating the lease
            $lease = [
                'user_id'  => Auth::User()->id,
                'group_id' => $group_id,
                'lease_ip' => $this->getClientIp().'/32',
                'lease_type' => 'aws',
                'protocol' => $protocol,
                'port_from'=> $port_from,
                'port_to'  => $port_to,
                'expiry'   => $expiry,
            ];

            $existingLease = Models\Lease::where('lease_ip', '=', $lease['lease_ip'])
                                    ->where('group_id', '=', $lease['group_id'])
                                    ->where('protocol', '=', $lease['protocol'])
                                    ->where('port_from', '=', $lease['port_from'])
                                    ->where('lease_type', '=', $lease['lease_type'])
                                    ->where('port_to', '=', $lease['port_to']);

            if ($existingLease->count() > 0) {
                $newLease = $existingLease->first();
                $newLease->expiry = $lease['expiry'];
                $newLease->save();
            } else {
                $result = $this->createLease($lease);
                if (! $result) {
                    //Lease Creation Failed. AWS Reported an error. Generally in case if a lease with same ip, protocol, port already exists on AWS.
                    return redirect("/manage/$group_id")
                                    ->with('message', 'Lease Creation Failed! Does a similar lease already exist? Terminate that first.');
                }
                $lease = Models\Lease::create($lease);
            }

            $this->NotificationMail($lease, true);

            return redirect("/manage/$group_id")
                        ->with('message', 'Lease created successfully!');
        } elseif (2 == $input['access']) {
            $email = $input['email'];
        } elseif (4 == $input['access']) {
            $email = 'DEPLOY';
        }

        $token = md5(time() + rand());
        $invite = [
            'user_id'  => Auth::User()->id,
            'group_id' => $group_id,
            'protocol' => $protocol,
            'port_from'=> $port_from,
            'port_to'  => $port_to,
            'expiry'   => $expiry,
            'email'    => $email,
            'token'    => $token,
        ];
        $invite = Models\Invite::create($invite);
        if ($email && $email != 'DEPLOY') {
            $data = ['invite'=>$invite->toArray()];
            //Send Invite Mail
            Mail::queue('emails.invite', $data, function ($message) use ($email) {
                $message->to($email, 'Invite')->subject('Access Lease Invite');
            });

            return redirect("/manage/$group_id")
                           ->with('message', 'Invite Sent successfully!');
        } else {
            return View::make('pages.invited')->with('invite', $invite);
        }
    }

    /*
      * Terminates the active leases & invites
      * @return getManage View
    */
    public function postTerminate($group_id)
    {
        $input = Request::all();
        if (isset($input['invite_id'])) {
            //Terminate Invite
            // Check for existence of invite
            try {
                $invite = Models\Invite::findorFail($input['invite_id']);
            } catch (Exception $e) {
                $message = 'Invite not found';

                return redirect("/manage/$group_id")->with('message', $message);
            }
            $invite->delete();

            return redirect("/manage/$group_id")
                                ->with('message', 'Invite terminated successfully');
        } elseif (isset($input['lease_id'])) {
            //Terminate Lease
            // Check for existence of lease
            try {
                $lease = Models\Lease::findorFail($input['lease_id']);
            } catch (Exception $e) {
                $message = 'Lease not found';

                return redirect("/manage/$group_id")->with('message', $message);
            }
            // Terminate the lease on AWS
            $result = $this->terminateLease($lease->toArray());
            //Delete from DB
            $lease->delete();
            $this->NotificationMail($lease, false);

            if (! $result) {
                //Should not occur even if lease doesn't exist with AWS. Check AWS API Conf.
                return redirect("/manage/$group_id")
                                ->with('message', 'Lease Termination returned error. Assumed the lease was already deleted');
            }

            return redirect("/manage/$group_id")
                                ->with('message', 'Lease terminated successfully');
        } else {
            App::abort(403, 'Unauthorized action.');
        }

        //get security group details
        $ec2 = App::make('aws')->createClient('ec2');
        $security_group = $ec2->describeSecurityGroups([
            'GroupIds' => [$group_id],
        ]);
        $security_group = $security_group['SecurityGroups'][0];

        //get Active Leases
        $leases = Models\Lease::getByGroupId($group_id);

        //get Active Invites
        $invites = Models\Invite::getByGroupId($group_id);

        return View::make('getManage')
                    ->with('security_group', $security_group)
                    ->with('leases', $leases)
                    ->with('invites', $invites);
    }

    /*
     * Handles cleaning of expired lease, called via artisan command custom:leasemanager run via cron
     * return void
     */

    public function cleanLeases()
    {
        $messages = [];
        $leases = Models\Lease::where('lease_type', 'aws')->get();
        foreach ($leases as $lease) {
            $time_left = strtotime($lease->created_at) + $lease->expiry - time();
            if ($time_left <= 0) {
                $result = $this->terminateLease($lease->toArray());
                $lease->delete();
                $this->NotificationMail($lease, false);
                if (! $result) {
                    array_push($messages, "Lease Termination of Lease ID $lease->id reported error on AWS API Call. Assumed already deleted.");
                }
            }
        }
        if (! empty($messages)) {
            return implode("\n", $messages);
        }

        return "No leases to clear";
    }

    /*
     * Handles Guest Access for lease invites
     */
    public function getInvite($token)
    {
        $invite = Models\Invite::getByToken($token);
        if (! $invite) {
            return View::make('pages.guest')->with('failure', 'Invalid Token. It was already used or has been terminated by the admins');
        }

        $email = $invite->email;
        if (! $invite->email) {
            $email = 'URL';
        }
        //Creating the lease
            $lease = [
                'user_id'     => $invite->user_id,
                'group_id'    => $invite->group_id,
                'lease_ip'    => $this->getClientIp().'/32',
                'protocol'    => $invite->protocol,
                'port_from'   => $invite->port_from,
                'port_to'     => $invite->port_to,
                'expiry'      => $invite->expiry,
                'invite_email'=> $email,
            ];

        $result = $this->createLease($lease);

        if (! $result) {
            //Lease Creation Failed. AWS Reported an error. Generally in case if a lease with same ip, protocl, port already exists on AWS.
                return View::make('pages.guest')->with('failure', "Error encountered while creating lease. Please try again. If doesn't help contact the admin.");
        }
        $lease = Models\Lease::create($lease);
        if ($invite->email != 'DEPLOY') {
            $invite = $invite->delete();
        }
        $this->NotificationMail($lease, true);

        return View::make('pages.guest')->with('lease', $lease);
    }

    /*
     * Handles Display of details of site users only to site admin
     */
    public function getUsers()
    {
        $users = Models\User::get();

        return view('getUsers', [
            'users' => $users,
        ]);
    }

    /*
     * Handles Display of new user form (only to site admin)
     */
    public function getAddUser()
    {
        $user = new Models\User();

        return View::make('getAddUser', compact('user'));
    }

    /*
     * Handles Adding of new user (only for site admin)
     */
    public function postAddUser()
    {
        $input = Request::all();

        //Validation Rules
        $user_rules = [
            'email'    => 'required|between:2,50|email|unique:users|org_email',
            'name'     => 'required|between:3,100',
            'admin'    => 'required|in:1,0',
        ];

        $domain = config('concierge.google_domain');

        $validator = Validator::make($input, $user_rules, [
            'org_email' => "Only $domain emails allowed",
        ]);

        if ($validator->fails()) {
            return redirect('/users/add')
                ->with('errors', $validator->messages()->toArray());
        } else {

            // Backward compatible
            $input['password'] = '';
            Models\User::create($input);

            return redirect('/users')
                            ->with('message', 'User Added Successfully');
        }
    }

    public function getEditUser($id)
    {
        $user = Models\User::find($id);

        return View::make('getAddUser', compact('user'));
    }

    public function postEditUser($id)
    {
        $input = Request::all();

        //Validation Rules
        $user_rules = [
            'email'              => "required|between:2,50|email|unique:users,email,$id|org_email",
            'name'               => 'required|between:3,100|alpha_spaces',
            'admin'              => 'required|in:1,0',
        ];

        $validator = Validator::make($input, $user_rules, [
            'org_email' => 'Only razorpay.com emails allowed',
        ]);

        if ($validator->fails()) {
            return redirect("/user/$id/edit")->with('errors', $validator->messages()->toArray());
        } else {
            Models\User::find($id)->update($input);

            return redirect('/users')->with('message', 'User Saved Successfully');
        }
    }

    /*
     * Handles deletion of users (only for site admin)
     */
    public function postUsers()
    {
        $input = Request::all();
        $message = null;

        if (! isset($input['user_id'])) {
            App::abort(403, 'Unauthorized action.');
        }
        try {
            $user = Models\User::findorfail($input['user_id']);
        } catch (Exception $e) {
            //User not found
            return redirect('/users')
                    ->with('message', 'Invalid User');
        }

        if ($user->id == Auth::user()->id) {
            //Avoid Self Delete
            $message = "You can't delete yourself";
        } else {
            $deleted = $user->delete();
            $message = 'User Deleted Successfully';
        }

        return redirect('/users')
                    ->with('message', $message);
    }

    /*
     * Handles sending of notification mail
     * Requires two arguements $lease, $ mode.
     * $lease = Lease Object Containing the lease created or deleted
     * $mode = TRUE for lease created, FALSE for lease deleted
     */

    private function NotificationMail($lease, $mode)
    {
        $data = ['lease'=>$lease->toArray(), 'mode'=>$mode];

        if ($mode) {
            //In case of Lease Creation
            $username = $lease->user->username;
            $type = (isset($lease['invite_email'])) ? (('URL' == $lease['invite_email']) ? 'URL Invite' : $lease['invite_email']) : 'Self';

            Log::info('Secure Lease Created at: '.$lease['created_at'].", Creator: $username, Type: $type, Group: ".$lease['group_id'].
                ', Leased IP: '.$lease['lease_ip'].', Ports: '.$lease['port_from'].
                '-'.$lease['port_to'].', Protocol: '.$lease['protocol'].', Expiry: '.$lease['expiry']);

            Mail::queue('emails.notification', $data, function ($message) {
                $message->to(config('concierge.notification_emailid'), 'Security Notification')->subject('Secure Access Lease Created');
            });
        } else {
            //In Case of Lease Termination
            $username = $lease->user->username;
            $type = (isset($lease['invite_email'])) ? (('URL' == $lease['invite_email']) ? 'URL Invite' : $lease['invite_email']) : 'Self';
            $terminator = (null !== Auth::user()) ? Auth::user()->username : 'Self-Expiry';

            Log::info('Secure Lease Terminated at: '.$lease['deleted_at'].", Creator: $username, Type: $type, Group: ".$lease['group_id'].
                ', Leased IP: '.$lease['lease_ip'].', Ports: '.$lease['port_from'].
                '-'.$lease['port_to'].', Protocol: '.$lease['protocol'].", Terminated By: $terminator");

            Mail::queue('emails.notification', $data, function ($message) {
                $message->to(config('concierge.notification_emailid'), 'Security Notification')->subject('Secure Access Lease Terminated');
            });
        }
    }

    /*
     * Handles lease creation by communitacting with AWS API
     * Requires an associative array of lease row.
     * return true if successful, false when AWS API returns error
     *
     * @see https://docs.aws.amazon.com/aws-sdk-php/v3/api/api-ec2-2016-11-15.html#updatesecuritygroupruledescriptionsingress
     */
    private function createLease($lease)
    {
        $ec2 = App::make('aws')->createClient('ec2');

        $email = Auth::User()->email;
        $time = Carbon::now()->toDateTimeString();

        $permissions = [
            [
                'FromPort'   => $lease['port_from'],
                'IpProtocol' => $lease['protocol'],
                'ToPort'     => $lease['port_to'],
                'IpRanges'   => [
                    [
                        'CidrIp'        => $lease['lease_ip'],
                        'Description'   => "Created by $email at $time",
                    ]
                ],
            ]
        ];
        try {

            $result = $ec2->authorizeSecurityGroupIngress([
                'DryRun'     => false,
                'GroupId'    => $lease['group_id'],
                'IpPermissions' => $permissions
            ]);

        } catch (Exception $e) {

            Log::info('Error while creating lease', [
                'lease'     =>  $lease,
                'exception' =>  $e->getMessage(),
            ]);

            return false;
        }

        return true;
    }

    /*
     * Handles lease termination by communitacting with AWS API
     * Requires an associative array of lease row.
     * return true if successful, false when AWS API returns error
     */
    private function terminateLease($lease)
    {
        $ec2 = App::make('aws')->createClient('ec2');
        try {
            $result = $ec2->revokeSecurityGroupIngress([
            'DryRun'     => false,
            'GroupId'    => $lease['group_id'],
            'IpProtocol' => $lease['protocol'],
            'FromPort'   => $lease['port_from'],
            'ToPort'     => $lease['port_to'],
            'CidrIp'     => $lease['lease_ip'],
            ]);
        } catch (Exception $e) {
            return false;
        }

        return true;
    }

    private function getClientIp()
    {
        if (isset($_SERVER['HTTP_X_FORWARDED_FOR']) and $_SERVER['HTTP_X_FORWARDED_FOR']) {
            // if behind an load balancer, assume that all load balancers have private IPs
            // and the first public IP will be that of the client
            $clientIpAddress = $_SERVER['HTTP_X_FORWARDED_FOR'];

            // We pick the first non-public IP we get

            if (strpos($clientIpAddress, ',') !== false)
            {
                $ips = array_reverse(array_map('trim', explode(',' , $clientIpAddress)));

                function isPublicIp($ip)
                {
                    $flags = FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE;

                    return (bool) filter_var(
                        $ip,
                        FILTER_VALIDATE_IP,
                        [
                            'flags' => $flags,
                        ]
                    );
                }

                foreach ($ips as $ip)
                {
                    if (isPublicIp($ip))
                    {
                        return $ip;
                    }
                }
            }
        } else {
            // if not behind ELB
            $clientIpAddress = $_SERVER['REMOTE_ADDR'];
        }

        return $clientIpAddress;
    }

    public function getStatus()
    {
        $msgArray = [];

        try
        {
            if (DB::connection('mysql')->getPdo())
            {
                $msgArray = [
                    'msg' => 'Connected to DB',
                ];
            }

            return Response::json($msgArray);
        }
        catch(\Exception $e)
        {
            return Response::json(['error' => 'DB Connection error'], 500);
        }
    }
}
