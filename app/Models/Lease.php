<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Lease extends Model
{
    use \Illuminate\Database\Eloquent\SoftDeletes;

    protected $guarded = ['id'];

    //Deleting only changes deleted_at column value, doesnt remove the field for tracking
    protected $softDelete = true;

    protected $dates = ['created_at', 'updated_at', 'deleted_at'];

    // Lease __belongs_to__ User
    public function user()
    {
        return $this->belongsTo('App\Models\User');
    }

    /*
     *Returns active leases by groupId
     */
    public static function getByGroupId($group_id)
    {
        $leases = self::where('group_id', $group_id)->get();

        return $leases;
    }
}
