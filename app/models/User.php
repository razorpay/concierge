<?php

use Illuminate\Auth\UserInterface;
use Illuminate\Auth\Reminders\RemindableInterface;

class User extends Eloquent implements UserInterface, RemindableInterface {

	/**
	 * The database table used by the model.
	 *
	 * @var string
	 */
	protected $table = 'users';

	protected $gaurded = array('id');

	//Editable fields
	protected $fillable = array('name', 'password');
	
	/**
	 * The attributes excluded from the model's JSON form.
	 *
	 * @var array
	 */
	protected $hidden = array('password');
	
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
	 * @param  string  $value
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

    public function leases()
    {
        return $this->hasMany('Lease');
    }

    public function invites()
    {
        return $this->hasMany('Invite');
        
    }

    public static function getIdFromUsername($username)
    {
        $query = DB::table('users')
                 ->where('username', '=', $username)
                 ->first();
        return $query->id;
    }

    public function getActiveInvites($type=NULL)
    {
        if(!$type)
        {
        	return $this->invites();
	    }
	    elseif("url"==$type)
	    {
	    	$invites = $this->invites->filter(function($invite)
			{
			    if(!$invite->email) return true;
			});
			return $invites;
	    }
	    elseif("email"==$type)
	    {
	    	$invites = $this->invites->filter(function($invite)
			{
			    if($invite->email) return true;
			});
			return $invites;
	    }

        
    }

    public function getActiveLeases($type=NULL)
    {
        if(!$type)
        {
        	return $this->leases();
	    }
	    elseif("self"==$type)
	    {
	    	$leases = $this->leases->filter(function($lease)
			{
			    if(! $lease->invite_email) return true;
			});
			return $leases;
	    }
	    elseif("url"==$type)
	    {
	    	$leases = $this->leases->filter(function($lease)
			{
			    if("NoEmail"==$lease->invite_email) return true;
			});
			return $leases;
	    }
	    elseif("email"==$type)
	    {
	    	$leases = $this->leases->filter(function($lease)
			{
			    if("NoEmail"!=$lease->invite_email && $lease->invite_email) return true;
			});
			return $leases;
	    }

        
    }

    public function delete()
    {
    	Lease::where("user_id", $this->id)->delete();
    	Invite::where("user_id", $this->id)->delete();
    	return parent::delete();
    }
}
