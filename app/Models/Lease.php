<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\SoftDeletes;
use AWS;
use Carbon\Carbon;
use Exception;
use Illuminate\Support\Facades\Log;

class Lease extends Model
{
    use SoftDeletes;

    /**
     * The database table used by the model.
     *
     * @var string
     */
    protected $table = 'leases';

    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'user_id', 'group_id', 'lease_ip', 'protocol', 'port_from', 'port_to',
        'expiry', 'invite_email',
    ];

    /**
     * Lease __belongs_to__ User
     *
     * @return void
     */
    public function user()
    {
        return $this->belongsTo('App\Models\User');
    }

    /**
     * Returns active leases by groupId
     *
     * @param  int $group_id
     * @return void
     */
    public static function getByGroupId($group_id)
    {
        $leases = self::where('group_id', $group_id)->get();

        return $leases;
    }

    /**
     * Handles lease creation by communitacting with AWS API
     * Requires an associative array of lease row.
     * return true if successful, false when AWS API returns error
     *
     * @param array $lease
     * @param string $email
     * @return boolean
     */
    public static function createLease($lease, $email)
    {
        $ec2 = AWS::createClient('ec2');
        $time = Carbon::now()->toDateTimeString();

        $permissions = [
            [
                'FromPort'   => $lease['port_from'],
                'IpProtocol' => $lease['protocol'],
                'ToPort'     => $lease['port_to'],
                'IpRanges'   => [
                    [
                        'CidrIp'        => $lease['lease_ip'],
                        'Description'   => "Created by $email at $time",
                    ]
                ],
            ]
        ];
        try {
            $ec2->authorizeSecurityGroupIngress([
                'DryRun'     => false,
                'GroupId'    => $lease['group_id'],
                'IpPermissions' => $permissions
            ]);
        } catch (Exception $e) {
            Log::info('Error while creating lease', [
                'lease'     =>  $lease,
                'exception' =>  $e->getMessage(),
            ]);
            return false;
        }
        return true;
    }

    /**
     * Handles lease termination by communitacting with AWS API
     * Requires an associative array of lease row.
     * return true if successful, false when AWS API returns error
     *
     * @param array $lease
     * @return boolean
     */
    public static function terminateLease($lease)
    {
        $ec2 = AWS::createClient('ec2');
        try {
            $ec2->revokeSecurityGroupIngress([
            'DryRun'     => false,
            'GroupId'    => $lease['group_id'],
            'IpProtocol' => $lease['protocol'],
            'FromPort'   => $lease['port_from'],
            'ToPort'     => $lease['port_to'],
            'CidrIp'     => $lease['lease_ip'],
            ]);
        } catch (Exception $e) {
            Log::info('Error while terminating lease', [
                'lease'     =>  $lease,
                'exception' =>  $e->getMessage(),
            ]);
            return false;
        }
        return true;
    }

    public static function cleanLeases()
    {
        $messages = [];
        $leases = self::get();
        foreach ($leases as $lease) {
            $time_left = strtotime($lease->created_at) + $lease->expiry - time();
            if ($time_left <= 0) {
                $result = self::terminateLease($lease->toArray());
                $lease->delete();
                if (! $result) {
                    array_push($messages, "Lease Termination of Lease ID $lease->id reported error on AWS API Call. Assumed already deleted.");
                }
            }
        }
        if (! empty($messages)) {
            return implode("\n", $messages);
        }
        return "No leases to clear";
    }
}
