<?php

namespace App\Models;

use Illuminate\Auth\Authenticatable;
use Illuminate\Auth\Passwords\CanResetPassword;
use Illuminate\Contracts\Auth\Authenticatable as AuthenticatableContract;
use Illuminate\Contracts\Auth\CanResetPassword as CanResetPasswordContract;
use Illuminate\Database\Eloquent\Model;

class User extends Model implements AuthenticatableContract, CanResetPasswordContract
{
    use Authenticatable, CanResetPassword;
    use \Illuminate\Database\Eloquent\SoftDeletes;

    /**
     * The database table used by the model.
     *
     * @var string
     */
    protected $table = 'users';

    protected $gaurded = ['id'];

    //Editable fields
    protected $fillable = ['google_id', 'email', 'access_token', 'username', 'name', 'password', 'admin'];

    protected $dates = ['created_at', 'updated_at', 'deleted_at'];

    /**
     * The attributes excluded from the model's JSON form.
     *
     * @var array
     */
    protected $hidden = ['password'];

    /**
     * Get the unique identifier for the user.
     *
     * @return mixed
     */
    public function getAuthIdentifier()
    {
        return $this->getKey();
    }

    /**
     * Get the password for the user.
     *
     * @return string
     */
    public function getAuthPassword()
    {
        return $this->password;
    }

    /**
     * Get the token value for the "remember me" session.
     *
     * @return string
     */
    public function getRememberToken()
    {
        return $this->remember_token;
    }

    /**
     * Set the token value for the "remember me" session.
     *
     * @param string $value
     *
     * @return void
     */
    public function setRememberToken($value)
    {
        $this->remember_token = $value;
    }

    /**
     * Get the column name for the "remember me" token.
     *
     * @return string
     */
    public function getRememberTokenName()
    {
        return 'remember_token';
    }

    /**
     * Get the e-mail address where password reminders are sent.
     *
     * @return string
     */
    public function getReminderEmail()
    {
        return $this->email;
    }

    /**
     * One to Many relation with Lease.
     *
     * @return collection
     */
    public function leases()
    {
        return $this->hasMany('App\Models\Lease');
    }

    /**
     * One to Many relation with Invite.
     *
     * @return collection
     */
    public function invites()
    {
        return $this->hasMany('App\Models\Invite');
    }

    /**
     * Returns the id of a username.
     *
     * @return int
     */
    public static function getIdFromUsername($username)
    {
        $query = DB::table('users')
                 ->where('username', '=', $username)
                 ->first();

        return $query->id;
    }

    /**
     * Returns a list of active invites by type
     * Type can be NULL/Unspecified for all invites, "url" for URL Invites, "email" for email invites.
     *
     * @return collection
     */
    public function getActiveInvites($type = null)
    {
        if (!$type) {
            return $this->invites();
        } elseif ('url' == $type) {
            $invites = $this->invites->filter(function ($invite) {
                if (!$invite->email) {
                    return true;
                }
            });

            return $invites;
        } elseif ('deploy' == $type) {
            $invites = $this->invites->filter(function ($invite) {
                if ($invite->email == 'DEPLOY') {
                    return true;
                }
            });

            return $invites;
        } elseif ('email' == $type) {
            $invites = $this->invites->filter(function ($invite) {
                if ($invite->email && $invite->email != 'DEPLOY') {
                    return true;
                }
            });

            return $invites;
        }
    }

    /**
     * Returns a list of active leases by type
     * Type can be NULL/Unspecified for all leases, "self" for self leases, "url" for URL leases, "email" for email leases.
     *
     * @return collection
     */
    public function getActiveLeases($type = null)
    {
        if (!$type) {
            return $this->leases();
        } elseif ('self' == $type) {
            $leases = $this->leases->filter(function ($lease) {
                if (!$lease->invite_email) {
                    return true;
                }
            });

            return $leases;
        } elseif ('url' == $type) {
            $leases = $this->leases->filter(function ($lease) {
                if ('URL' == $lease->invite_email) {
                    return true;
                }
            });

            return $leases;
        } elseif ('email' == $type) {
            $leases = $this->leases->filter(function ($lease) {
                if ('URL' != $lease->invite_email && $lease->invite_email) {
                    return true;
                }
            });

            return $leases;
        }
    }
}
