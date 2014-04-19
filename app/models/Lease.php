<?php

class Lease extends Eloquent
{
	protected $guarded = array('id');
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