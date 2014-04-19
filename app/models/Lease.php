<?php

class Lease extends Eloquent
{
	protected $guarded = array('id');
    // Lease __belongs_to__ User
    public function user()
    {
        return $this->belongsTo('User');
    }

}