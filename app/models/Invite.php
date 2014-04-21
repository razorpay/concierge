<?php

class Invite extends Eloquent
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
     *Returns invite from token 
     */
    public static function getByToken($token)
	{
		try{
            $invite = self::where('token', $token)->firstOrFail();
        }
        catch(Exception $e)
        {
            $invite = FALSE;
        }
		return $invite;
	}
}