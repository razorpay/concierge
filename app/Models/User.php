<?php

namespace App\Models;

use Illuminate\Contracts\Auth\MustVerifyEmail;
use Illuminate\Foundation\Auth\User as Authenticatable;
use Illuminate\Notifications\Notifiable;
use Illuminate\Database\Eloquent\SoftDeletes;
use DB;

class User extends Authenticatable
{
    use Notifiable, SoftDeletes;

    /**
     * The database table used by the model.
     *
     * @var string
     */
    protected $table = 'users';

    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'name', 'email', 'password', 'google_id', 'access_token',
        'username', 'admin',
    ];

    /**
     * The attributes that should be hidden for arrays.
     *
     * @var array
     */
    protected $hidden = [
        'password', 'remember_token',
    ];

    /**
     * The attributes that should be cast to native types.
     *
     * @var array
     */
    protected $casts = [
        'email_verified_at' => 'datetime',
    ];

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
     * Returns a list of active leases by type
     *
     * @return collection
     */
    public function getActiveLeases()
    {
        return $this->leases();
    }
}
