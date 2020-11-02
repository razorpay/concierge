<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use SocialOAuth;
use App\Models\User;
use Illuminate\Support\Facades\Auth;
use App;
use Illuminate\Support\Facades\View;
use Illuminate\Support\Facades\Validator;
use Exception;
use App\Rules\ValidateOrganizationEmail;
use Carbon\Carbon;

class UserController extends Controller
{
    /**
     * Stage One - The Login form.
     *
     * @return Login View
     */
    public function getIndex(Request $request)
    {
        $code = $request->get('code');
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
            if ($result->verified_email !== true or
                (isset($result->hd) and $result->hd !== config('concierge.google_domain'))) {
                return App::abort(404);
            }

            // Find the user by email
            $user = User::where('email', $result->email)->first();
            if ($user) {
                // Update some fields
                $user->access_token = $token->getAccessToken();
                $user->google_id = $result->id;
                $user->save();

                // Login the user into the app
                Auth::loginUsingId($user->id);
                return redirect('/groups');
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

    /*
     * Handles Display of details of site users only to site admin
     */
    public function getUsers()
    {
        if (!Auth::User()->admin) {
            return redirect()->back()
            ->with('message', 'You don\'t have permission to manage users')
            ->with('class', 'Warning');
        }
        $users = User::get();

        return view('getUsers', [
            'users' => $users,
        ]);
    }

    /*
     * Handles Display of new user form (only to site admin)
     */
    public function getAddUser()
    {
        if (!Auth::User()->admin) {
            return redirect()->back()
            ->with('message', 'You don\'t have permission to add new users')
            ->with('class', 'Warning');
        }
        $user = new User();

        return View::make('getAddUser', compact('user'));
    }

    /*
     * Handles Adding of new user (only for site admin)
     */
    public function postAddUser(Request $request)
    {
        if (!Auth::User()->admin) {
            return redirect()->back()
            ->with('message', 'You don\'t have permission to add new users')
            ->with('class', 'Warning');
        }
        $input = $request->all();
        $input['username'] = explode('@', $input['email'])[0];

        $user_rules = [
            'email' => [
                "required",
                "between:2,50",
                "email",
                "unique:users,email,NULL,id,deleted_at,NULL",
                new ValidateOrganizationEmail()
            ],
            'username' => 'required|unique:users,username',
            'name'  => 'required|between:3,100',
            'admin' => 'required|in:1,0',
        ];

        $validator = Validator::make($input, $user_rules);

        if ($validator->fails()) {
            return redirect('/users/add')
                ->with('errors', $validator->messages());
        } else {
            User::create($input);

            return redirect('/users')
                ->with('message', 'User Added Successfully')
                ->with('class', 'Success');
        }
    }

    public function getEditUser($id)
    {
        if (!Auth::User()->admin) {
            return redirect()->back()
            ->with('message', 'You don\'t have permission to edit/update users')
            ->with('class', 'Warning');
        }
        $user = User::find($id);
        return View::make('getAddUser', compact('user'));
    }

    public function postEditUser(Request $request, $id)
    {
        if (!Auth::User()->admin) {
            return redirect()->back()
            ->with('message', 'You don\'t have permission to edit/update users')
            ->with('class', 'Warning');
        }
        $input = $request->all();
        $input['username'] = explode('@', $input['email'])[0];

        //Validation Rules
        $user_rules = [
            'email' => [
                "required",
                "between:2,50",
                "email",
                "unique:users,email," . $id . ",id,deleted_at,NULL",
                new ValidateOrganizationEmail()
            ],
            'username' => 'required|unique:users,username,' . $id,
            'name'  => 'required|between:3,100',
            'admin' => 'required|in:1,0',
        ];

        $validator = Validator::make($input, $user_rules);

        if ($validator->fails()) {
            return redirect("/user/$id/edit")->with('errors', $validator->messages());
        } else {
            User::find($id)->update($input);
            return redirect('/users')
                ->with('message', 'User Updated Successfully')
                ->with('class', 'Success');
        }
    }

    /*
     * Handles deletion of users (only for site admin)
     */
    public function postUsers(Request $request)
    {
        if (!Auth::User()->admin) {
            return redirect()->back()
            ->with('message', 'You don\'t have permission to delete users')
            ->with('class', 'Warning');
        }
        $input = $request->all();
        $message = null;

        if (! isset($input['user_id'])) {
            App::abort(403, 'Unauthorized action.');
        }
        try {
            $user = User::findorfail($input['user_id']);
        } catch (Exception $e) {
            return redirect('/users')
                ->with('message', 'Invalid User')
                ->with('class', 'Warning');
        }
        if ($user->id == Auth::user()->id) {
            //Avoid Self Delete
            $message = "You can't delete yourself";
        } else {
            $user->username = $user->username . Carbon::now();
            $user->save();
            $user->delete();
            $message = 'User Deleted Successfully';
        }

        return redirect('/users')
            ->with('message', $message)
            ->with('class', 'Success');
    }
}
