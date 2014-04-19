<?php

class Lease extends Eloquent
{
    // Lease __belongs_to__ User
    public function user()
    {
        return $this->belongsTo('User');
    }

}