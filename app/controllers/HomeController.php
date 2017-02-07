<?php

use App\LaravelDuo\LaravelDuo;
use Artdarek\OAuth\Facade\OAuth;

class HomeController extends BaseController
{
    // 6 hours
    const MAX_EXPIRY = 21600;

	/*
	|--------------------------------------------------------------------------
	| Default Home Controller
	|--------------------------------------------------------------------------
	|
	| You may wish to use controllers instead of, or in addition to, Closure
	| based routes. That's great! Here is an example controller method to
	| get you started. To route to this controller, just add the route:
	|
	|	Route::get('/', 'HomeController@showWelcome');
	|
	*/
        private $_laravelDuo;

    function __construct(LaravelDuo $laravelDuo)
    {
        $this->_laravelDuo = $laravelDuo;
    }

    /**
     * Stage One - The Login form
     * @return  Login View
     */
    public function getIndex()
    {
		/*
			Google Auth Logic
		*/

		$code = Input::get('code');

		$google_service = OAuth::consumer('Google');

		if (!$code)
		{
			$url = $google_service->getAuthorizationUri();

			return Redirect::to((string) $url);
		}
		else
		{
			$token = $google_service->requestAccessToken($code);

			$response = $google_service->request(Config::get('oauth-4-laravel.userinfo_url'));
			$result = json_decode($response);

			// Email must be:
			// - Verified
			// - Belong to razorpay.com domain
			//
			// Then only we'll create a user entry in the system or check for one
			if (!$result->verified_email or !checkEmailDomain($result->email)) return App::abort(404);

			// Find the user by email
			$user = User::where('email', $result->email)->first();

			if ($user)
			{
				// Update some fields
				$user->access_token = $token->getAccessToken();
				$user->google_id = $result->id;
				$user->password = ''; // backward compatibility

				$user->save();

				// Login the user into the app
				Auth::loginUsingId($user->id);
			}

			// whether logged in or not, just redirect
			return Redirect::to('/');
		}

		// This won't be rendered anymore, will remove once
		// Google auth is working well on production
        return View::make('pages.login');
    }

    /**
     * Stage Two - The Duo Auth form
     * @return Duo Login View or Redirect on error
     */
    public function postSignin()
    {
        $user = array(
            'username' => Input::get('username'),
            'password' => Input::get('password')
        );

        /**
         * Validate the user details, but don't log the user in
         */
        if(Auth::validate($user))
        {
            $U    = Input::get('username');

            $duoinfo = array(
                'HOST' => $this->_laravelDuo->get_host(),
                'POST' => URL::to('/') . '/duologin',
                'USER' => $U,
                'SIG'  => $this->_laravelDuo->signRequest($this->_laravelDuo->get_ikey(), $this->_laravelDuo->get_skey(), $this->_laravelDuo->get_akey(), $U)
            );

            return View::make('pages.duologin')->with(compact('duoinfo'));
        }
        else
        {
            return Redirect::to('/')->with('message', 'Your username and/or password was incorrect')->withInput();
        }

    }

    /**
     * Stage Three - After Duo Auth Form
     * @return Redirect to home
     */
    public function postDuologin()
    {
        /**
         * Sent back from Duo
         */
        $response = $_POST['sig_response'];

        $U = $this->_laravelDuo->verifyResponse($this->_laravelDuo->get_ikey(), $this->_laravelDuo->get_skey(), $this->_laravelDuo->get_akey(), $response);

        /**
         * Duo response returns USER field from Stage Two
         */
        if($U){

            /**
             * Get the id of the authenticated user from their email address
             */
            $id = User::getIdFromUsername($U);

            /**
             * Log the user in by their ID
             */
            Auth::loginUsingId($id);

            /**
             * Check Auth worked, redirect to homepage if so
             */
            if(Auth::check())
            {
                return Redirect::to('/');
            }
        }

        /**
         * Otherwise, Auth failed, redirect to homepage with message
         */
        return Redirect::to('/')->with('message', 'Unable to authenticate you.');

    }

    /**
     * Log user out
     * @return Redirect to Home
     */
    public function getLogout()
    {
        Auth::logout();
        return Redirect::to('/');
    }

	/**
	 * Get the list of all security groups
     * @return getGroups view
	 */
	public function getGroups()
	{
        //Get All security groups
		$ec2 = App::make('aws')->get('ec2');
		$security_groups=$ec2->describeSecurityGroups();
		$security_groups=$security_groups['SecurityGroups'];

        //Get all active leases
        $leases= Lease::get();

        //get all active Invites
        $invites= Invite::get();

        return View::make('getGroups')
                    ->with('security_groups', $security_groups)
                    ->with('leases', $leases)
                    ->with('invites', $invites);
	}

    /*
     * Displays a security groups details with active leases & security rules.
     * @return getManage View
     */
	public function getManage($group_id)
	{

        //get security group details
		$ec2 = App::make('aws')->get('ec2');
		$security_group=$ec2->describeSecurityGroups(array(
			'GroupIds' => array($group_id),
        ));
		$security_group=$security_group['SecurityGroups'][0];

        //get Active Leases
        $leases= Lease::getByGroupId($group_id);

        //get Active Invites
        $invites= Invite::getByGroupId($group_id);

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
        $input=Input::all();
        $messages = array();
        $email = NULL;
        /*
         For Lease Creation
        */
        if("ssh" == $input["rule_type"])
        {
            $protocol="tcp";
            $port_from="22";
            $port_to="22";
        }
        elseif("https" == $input["rule_type"])
        {
            $protocol="tcp";
            $port_from="443";
            $port_to="443";
        }
        elseif("custom" == $input["rule_type"])
        {
            $protocol=$input['protocol'];
            $port_from=$input['port_from'];
            $port_to=$input['port_to'];

            //Validations
            if ($protocol !== "tcp" and $protocol !== "udp")
            {
                array_push($messages, "Invalid Protocol");
            }
            if (!is_numeric($port_from) or $port_from > 65535 or $port_from <= 0)
            {
                array_push($messages, "Invalid From port");
            }
            if (!is_numeric($port_to) or $port_to > 65535 or $port_to <= 0)
            {
                array_push($messages, "Invalid To port");
            }
            if ($port_from > $port_to)
            {
                array_push($messages, "From port Must be less than equal to To Port");
            }
        }
        else
        {
            App::abort(403, 'Unauthorized action.');
        }

        //Other validations
        $expiry = $input['expiry'];

        if (!is_numeric($expiry) or $expiry <= 0 || $expiry > self::MAX_EXPIRY)
        {
            array_push($messages, "Invalid Expiry Time");
        }

        if (!in_array($input['access'], array(1, 2, 3, 4)))
        {
            array_push($messages, "Invalid invite Email");
        }

        if (2 == $input['access'])
        {
            if(!isset($input['email']) or !filter_var($input['email'], FILTER_VALIDATE_EMAIL))
            {
                array_push($messages, "Invalid invite Email");
            }
        }

        //Validation fails
        if(!empty($messages))
        {
            return Redirect::to("/manage/$group_id")
                            ->with('message', implode("<br/>", $messages));
        }

        if(1 == $input['access'])
        {
            //Creating the lease
            $lease = array(
                'user_id'=>Auth::User()->id,
                'group_id'=>$group_id,
                'lease_ip'=>$this->getClientIp()."/32",
                'protocol'=>$protocol,
                'port_from'=>$port_from,
                'port_to'=>$port_to,
                'expiry'=>$expiry,
            );

            $existingLease = Lease::where('lease_ip', '=', $lease['lease_ip'])
                                    ->where('group_id', '=', $lease['group_id'])
                                    ->where('protocol', '=', $lease['protocol'])
                                    ->where('port_from', '=', $lease['port_from'])
                                    ->where('port_to', '=', $lease['port_to']);

            if ($existingLease->count() > 0)
            {
                $newLease = $existingLease->first();
                $newLease->expiry = $lease['expiry'];
                $newLease->save();
            }
            else
            {
                $result=$this->createLease($lease);
                if(!$result)
                {
                    //Lease Creation Failed. AWS Reported an error. Generally in case if a lease with same ip, protocol, port already exists on AWS.
                    return Redirect::to("/manage/$group_id")
                                    ->with('message', "Lease Creation Failed! Does a similar lease already exist? Terminate that first.");
                }
                $lease=Lease::create($lease);
            }

            $this->NotificationMail($lease, TRUE);
            return Redirect::to("/manage/$group_id")
                        ->with('message', "Lease created successfully!");
        }
        elseif(2 == $input['access'])
        {
            $email=$input['email'];
        }
        elseif(4 == $input['access'])
        {
            $email='DEPLOY';
        }

        // TODO: Improve this
        $token=md5(time()+rand());
        $invite=array(
            'user_id'=>Auth::User()->id,
            'group_id'=>$group_id,
            'protocol'=>$protocol,
            'port_from'=>$port_from,
            'port_to'=>$port_to,
            'expiry'=>$expiry,
            'email'=>$email,
            'token'=>$token
        );

        $invite=Invite::create($invite);

        if($email and $email!='DEPLOY')
        {
            $data=array('invite'=>$invite->toArray());
            //Send Invite Mail
            Mail::queue('emails.invite', $data, function($message) use($email)
            {
                $message->to($email, 'Invite' )->subject('Access Lease Invite');
            });
            return Redirect::to("/manage/$group_id")
                           ->with('message', "Invite Sent successfully!");

        }
        else
        {
            return View::make('pages.invited')->with('invite', $invite);
        }
    }

    public function postRenew($group_id)
    {
        // $input = Input::all();
        // if (isset($input['invite_id']))
        // {
        //     try
        //     {
        //         $invite = Invite::findorfail($input['invite_id']);
        //     } catch (Exception $e) {
        //         $message="Invite not found";
        //         return Redirect::to("/manage/$group_id")->with('message', $message);
        //     }
        // }
        // else if (isset($input['lease_id']))
        // {
        //     // Check for existence of lease
        //     try
        //     {
        //         $lease=Lease::findorFail($input['lease_id']);
        //     }
        //     catch(Exception $e)
        //     {
        //         $message="Lease not found";
        //         return Redirect::to("/manage/$group_id")->with('message', $message);
        //     }
        //     // Terminate the lease on AWS
        //     $result=$this->terminateLease($lease->toArray());
        //     //Delete from DB
        //     $lease->delete();
        //     $lease = $lease->toArray();
        //     $this->NotificationMail($lease, FALSE);

        //     if(!$result)
        //     {
        //         //Should not occur even if lease doesn't exist with AWS. Check AWS API Conf.
        //         return Redirect::to("/manage/$group_id")
        //                         ->with('message', "Lease Termination returned error. Assumed the lease was already deleted");
        //     }

        //     $result=$this->createLease($lease);
        //     if(!$result)
        //     {
        //         //Lease Creation Failed. AWS Reported an error. Generally in case if a lease with same ip, protocol, port already exists on AWS.
        //         return Redirect::to("/manage/$group_id")
        //                         ->with('message', "Lease Creation Failed! Does a similar lease already exist? Terminate that first.");
        //     }
        //     $lease=Lease::create($lease);
        //     $this->NotificationMail($lease, TRUE);
        //     return Redirect::to("/manage/$group_id")
        //                 ->with('message', "Lease renewed successfully!");
        // }
    }

     /*
     * Terminates the active leases & invites
     * @return getManage View
     */
    public function postTerminate($group_id)
    {
        $input=Input::all();
        if(isset($input['invite_id']))
        {
            //Terminate Invite
            // Check for existence of invite
            try
            {
                $invite=Invite::findorFail($input['invite_id']);
            }
            catch(Exception $e)
            {
                $message="Invite not found";
                return Redirect::to("/manage/$group_id")->with('message', $message);
            }
            $invite->delete();
            return Redirect::to("/manage/$group_id")
                                ->with('message', "Invite terminated successfully");
        }
        elseif(isset($input['lease_id']))
        {
            //Terminate Lease
            // Check for existence of lease
            try
            {
                $lease=Lease::findorFail($input['lease_id']);
            }
            catch(Exception $e)
            {
                $message="Lease not found";
                return Redirect::to("/manage/$group_id")->with('message', $message);
            }
            // Terminate the lease on AWS
            $result=$this->terminateLease($lease->toArray());
            //Delete from DB
            $lease->delete();
            $this->NotificationMail($lease, FALSE);

            if(!$result)
            {
                //Should not occur even if lease doesn't exist with AWS. Check AWS API Conf.
                return Redirect::to("/manage/$group_id")
                                ->with('message', "Lease Termination returned error. Assumed the lease was already deleted");
            }
            return Redirect::to("/manage/$group_id")
                                ->with('message', "Lease terminated successfully");
        }
        else
        {
            App::abort(403, 'Unauthorized action.');
        }

        //get security group details
        $ec2 = App::make('aws')->get('ec2');
        $security_group=$ec2->describeSecurityGroups(array(
            'GroupIds' => array($group_id),
        ));
        $security_group=$security_group['SecurityGroups'][0];

        //get Active Leases
        $leases= Lease::getByGroupId($group_id);

        //get Active Invites
        $invites= Invite::getByGroupId($group_id);

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
        $messages=array();
        $leases=Lease::get();
        foreach($leases as $lease)
        {
            $time_left=strtotime($lease->created_at)+$lease->expiry-time();
            if($time_left<=0){
                $result=$this->terminateLease($lease->toArray());
                $lease->delete();
                $this->NotificationMail($lease, FALSE);
                if(!$result)
                {
                    array_push($messages,"Lease Termination of Lease ID $lease->id reported error on AWS API Call. Assumed already deleted.");
                }
            }
        }
        if(!empty($messages)) return implode("\n", $messages);
        return;
    }

    /*
     * Handles Guest Access for lease invites
     */
    public function getInvite($token)
    {
        $invite=Invite::getByToken($token);
        if(!$invite) return View::make('pages.guest')->with('failure', "Invalid Token. It was already used or has been terminated by the admins");

        $email=$invite->email;
        if(!$invite->email) $email="URL";
        //Creating the lease
            $lease=array(
                'user_id'=>$invite->user_id,
                'group_id'=>$invite->group_id,
                'lease_ip'=>$this->getClientIp()."/32",
                'protocol'=>$invite->protocol,
                'port_from'=>$invite->port_from,
                'port_to'=>$invite->port_to,
                'expiry'=>$invite->expiry,
                'invite_email'=>$email
            );
            $result=$this->createLease($lease);
            if(!$result)
            {
                //Lease Creation Failed. AWS Reported an error. Generally in case if a lease with same ip, protocl, port already exists on AWS.
                return View::make('pages.guest')->with('failure', "Error encountered while creating lease. Please try again. If doesn't help contact the admin.");
            }
            $lease=Lease::create($lease);
            if($invite->email != 'DEPLOY') $invite=$invite->delete();
            $this->NotificationMail($lease, TRUE);
            return View::make('pages.guest')->with('lease', $lease);
    }

    /*
     * Handles Display of details of site users only to site admin
     */
    public function getUsers()
    {
        $users=User::get();
        return View::make('getUsers')
                    ->with('users', $users);
    }

    /*
     * Handles Display of new user form (only to site admin)
     */
    public function getAddUser()
    {
		$user = new User;

        return View::make('getAddUser', compact('user'));
    }

    /*
     * Handles Adding of new user (only for site admin)
     */
    public function postAddUser()
    {
        $input = Input::all();
        //Validation Rules
        $user_rules = array(
        'email'              => 'required|between:2,50|email|unique:users|razorpay_email',
        'name'                  => 'required|between:3,100|alpha_spaces',
        'admin'                 => 'required|in:1,0');

        $validator = Validator::make($input, $user_rules, array(
			'razorpay_email' => 'Only razorpay.com emails allowed'
		));
        if ($validator->fails())
        {
             return Redirect::to('/users/add')
                            ->with('errors', $validator->messages()->toArray() );
        }
        else
        {
			$input['password'] = ''; // Backward compatible
            User::create($input);

            return Redirect::to('/users')
                            ->with('message', "User Added Successfully" );
        }

    }

	public function getEditUser($id)
	{
		$user = User::find($id);

		return View::make('getAddUser', compact('user'));
	}

	public function postEditUser($id)
	{

		$input = Input::all();

		//Validation Rules
        $user_rules = array(
			'email'              => "required|between:2,50|email|unique:users,email,$id|razorpay_email",
			'name'               => 'required|between:3,100|alpha_spaces',
			'admin'              => 'required|in:1,0'
		);

        $validator = Validator::make($input, $user_rules, array(
			'razorpay_email' => 'Only razorpay.com emails allowed'
		));

        if ($validator->fails())
        {
			return Redirect::to("/user/$id/edit")->with('errors', $validator->messages()->toArray());
        }
        else
        {
            User::find($id)->update($input);

            return Redirect::to('/users')->with('message', "User Saved Successfully" );
        }
	}


    /*
     * Handles deletion of users (only for site admin)
     */
    public function postUsers()
    {
        $input=Input::all();
        $message=NULL;

        if(!isset($input['user_id'])) App::abort(403, 'Unauthorized action.');
        try
        {
            $user=User::findorfail($input['user_id']);
        }
        catch(Exception $e)
        {
            //User not found
            return Redirect::to('/users')
                    ->with('message', "Invalid User");
        }

       if($user->id == Auth::user()->id)
        {
            //Avoid Self Delete
            $message="You can't delete yourself";
        }
        else
        {
            $deleted=$user->delete();
            $message="User Deleted Successfully";
        }
        return Redirect::to('/users')
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
        $data=array('lease'=>$lease->toArray(), 'mode'=>$mode);

        if($mode)
        {
            //In case of Lease Creation
            $username=$lease->user->username;
            $type=(isset($lease['invite_email'])) ? (("URL"==$lease['invite_email'])?"URL Invite":$lease['invite_email']) : "Self";


            Log::info("Secure Lease Created at: ".$lease['created_at'].", Creator: $username, Type: $type, Group: ".$lease['group_id'].
                ", Leased IP: ".$lease['lease_ip'].", Ports: ".$lease['port_from'].
                "-".$lease['port_to'].", Protocol: ".$lease['protocol'].", Expiry: ".$lease['expiry']);


            Mail::queue('emails.notification', $data, function($message)
            {
                $message->to(Config::get('custom_config.notification_emailid'), 'Security Notification' )->subject('Secure Access Lease Created');
            });
        }
        else
        {
            //In Case of Lease Termination
            $username=$lease->user->username;
            $type=(isset($lease['invite_email'])) ? (("URL"==$lease['invite_email'])?"URL Invite":$lease['invite_email']) : "Self";
            $terminator=(null !== Auth::user()) ? Auth::user()->username : "Self-Expiry";

            Log::info("Secure Lease Terminated at: ".$lease['deleted_at'].", Creator: $username, Type: $type, Group: ".$lease['group_id'].
                ", Leased IP: ".$lease['lease_ip'].", Ports: ".$lease['port_from'].
                "-".$lease['port_to'].", Protocol: ".$lease['protocol'].", Terminated By: $terminator");

            Mail::queue('emails.notification', $data, function($message)
            {
                $message->to(Config::get('custom_config.notification_emailid'), 'Security Notification' )->subject('Secure Access Lease Terminated');
            });
        }
    }

    /*
     * Handles lease creation by communitacting with AWS API
     * Requires an associative array of lease row.
     * return true if successful, false when AWS API returns error
     */
    private function createLease($lease)
    {
        $ec2 = App::make('aws')->get('ec2');
        try
        {
            $result = $ec2->authorizeSecurityGroupIngress(array(
            'DryRun' => false,
            'GroupId' =>  $lease['group_id'],
            'IpProtocol' => $lease['protocol'],
            'FromPort' => $lease['port_from'],
            'ToPort' => $lease['port_to'],
            'CidrIp' => $lease['lease_ip'],
            ));
        }
        catch(Exception $e)
        {
            return FALSE;
        }
        return TRUE;
    }
    /*
     * Handles lease termination by communitacting with AWS API
     * Requires an associative array of lease row.
     * return true if successful, false when AWS API returns error
     */
    private function terminateLease($lease)
    {
        $ec2 = App::make('aws')->get('ec2');
        try
        {
            $result = $ec2->revokeSecurityGroupIngress(array(
            'DryRun' => false,
            'GroupId' => $lease['group_id'],
            'IpProtocol' => $lease['protocol'],
            'FromPort' => $lease['port_from'],
            'ToPort' => $lease['port_to'],
            'CidrIp' => $lease['lease_ip'],
            ));
        }
        catch(Exception $e)
        {
            return FALSE;
        }
        return TRUE;
    }

    private function getClientIp()
    {
        if (isset($_SERVER['HTTP_X_FORWARDED_FOR']) and $_SERVER['HTTP_X_FORWARDED_FOR']) {
            // if behind an ELB
            $clientIpAddress = $_SERVER['HTTP_X_FORWARDED_FOR'];
        } else {
            // if not behind ELB
            $clientIpAddress = $_SERVER['REMOTE_ADDR'];
        }

        return $clientIpAddress;
    }

}
