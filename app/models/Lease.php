<?php

class Lease extends Eloquent
{
	protected $guarded = array('id');
	//Deleting only changes deleted_at column value, doesnt remove the field for tracking
	protected $softDelete = true;
    // Lease __belongs_to__ User
    public function user()
    {
        return $this->belongsTo('User');
    }

    /*
     *Returns active leases
     */
    public static function getByGroupId($group_id)
	{
		$leases = self::where('group_id', $group_id)->get();
		return $leases;
	}
}