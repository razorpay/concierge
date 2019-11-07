<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Models\Lease;
use Auth;
use App;
use View;
use Carbon\Carbon;
use Exception;
use Illuminate\Support\Facades\Log;
use AWS;

class GroupController extends Controller
{

    // 6 hours
    const MAX_EXPIRY = 21600;

    const CONCIERGE_TAG = 'concierge';

    /**
     * Get the list of all security groups.
     *
     * @return getGroups view
     */
    public function getGroups()
    {
        //Get All security groups
        $ec2 = AWS::createClient('ec2');

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

        //Get all active leases
        $leases = Lease::get();

        return view('getGroups', [
            'security_groups'   => $security_groups,
            'leases'            => $leases,
        ]);
    }

     /*
     * Displays a security groups details with active leases & security rules.
     * @return getManage View
     */
    public function getManage($group_id)
    {
        //get security group details
        $ec2 = AWS::createClient('ec2');
        $security_group = $ec2->describeSecurityGroups([
            'GroupIds' => [$group_id],
        ]);
        $security_group = $security_group['SecurityGroups'][0];

        //get Active Leases
        $leases = Lease::getByGroupId($group_id);

        return View::make('getManage')
            ->with('security_group', $security_group)
            ->with('leases', $leases);
    }

    /*
     * Handles Lease creation & termination post requests to Group Manage page
     * @return Redirect to getManage View with error/success
     */
    public function postManage(Request $request, $group_id)
    {
        $input = $request->all();
        $messages = [];
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

        //Validation fails
        if (! empty($messages)) {
            return redirect("/manage/$group_id")
                ->with('message', implode('<br/>', $messages))
                ->with('class', 'Danger');
        }

        //Creating the lease
        $lease = [
            'user_id'  => Auth::User()->id,
            'group_id' => $group_id,
            'lease_ip' => $this->getClientIp().'/32',
            'protocol' => $protocol,
            'port_from'=> $port_from,
            'port_to'  => $port_to,
            'expiry'   => $expiry,
        ];

        $existingLease = Lease::where('lease_ip', '=', $lease['lease_ip'])
            ->where('group_id', '=', $lease['group_id'])
            ->where('protocol', '=', $lease['protocol'])
            ->where('port_from', '=', $lease['port_from'])
            ->where('port_to', '=', $lease['port_to']);

        if ($existingLease->count() > 0) {
            // dd($existingLease);
            $newLease = $existingLease->first();
            $newLease->expiry = $lease['expiry'];
            $newLease->save();
        } else {
            $result = Lease::createLease($lease, Auth::User()->email);
            if (! $result) {
                // Lease Creation Failed.
                // AWS Reported an error. Generally in case if a lease with same ip,
                // protocol, port already exists on AWS.
                return redirect("/manage/$group_id")
                    ->with(
                        'message',
                        'Lease Creation Failed! Does a similar lease already exist? Terminate that first.'
                    )
                    ->with('class', 'Danger');
            }
            $lease = Lease::create($lease);
        }

        // $this->NotificationMail($lease, true);

        return redirect("/manage/$group_id")
            ->with('message', 'Lease created successfully!')
            ->with('class', 'Success');
    }

    private function getClientIp()
    {
        if (isset($_SERVER['HTTP_X_FORWARDED_FOR']) and $_SERVER['HTTP_X_FORWARDED_FOR']) {
            // if behind an load balancer, assume that all load balancers have private IPs
            // and the first public IP will be that of the client
            $clientIpAddress = $_SERVER['HTTP_X_FORWARDED_FOR'];

            // We pick the first non-public IP we get

            if (strpos($clientIpAddress, ',') !== false) {
                $ips = array_reverse(array_map('trim', explode(',', $clientIpAddress)));

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

                foreach ($ips as $ip) {
                    if (isPublicIp($ip)) {
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

    /*
      * Terminates the active leases & invites
      * @return getManage View
    */
    public function postTerminate(Request $request, $group_id)
    {
        $input = $request->all();
        if (isset($input['lease_id'])) {
            //Terminate Lease
            // Check for existence of lease
            try {
                $lease = Lease::findorFail($input['lease_id']);
            } catch (Exception $e) {
                $message = 'Lease not found';
                return redirect("/manage/$group_id")
                    ->with('message', $message)
                    ->with('class', 'Warning');
            }
            // Terminate the lease on AWS
            $result = Lease::terminateLease($lease->toArray());
            //Delete from DB
            $lease->delete();

            if (! $result) {
                //Should not occur even if lease doesn't exist with AWS. Check AWS API Conf.
                return redirect("/manage/$group_id")
                    ->with('message', 'Lease Termination returned error. Assumed the lease was already deleted')
                    ->with('class', 'Warning');
            }

            return redirect("/manage/$group_id")
                ->with('message', 'Lease terminated successfully')
                ->with('class', 'Success');
        } else {
            App::abort(403, 'Unauthorized action.');
        }

        //get security group details
        $ec2 = AWS::createClient('ec2');
        $security_group = $ec2->describeSecurityGroups([
            'GroupIds' => [$group_id],
        ]);
        $security_group = $security_group['SecurityGroups'][0];

        //get Active Leases
        $leases = Lease::getByGroupId($group_id);

        //get Active Invites
        return View::make('getManage')
            ->with('security_group', $security_group)
            ->with('leases', $leases);
    }
}
