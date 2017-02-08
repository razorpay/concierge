<?php

return array(
	/*
    |--------------------------------------------------------------------------
    | Paths
    |--------------------------------------------------------------------------
    |
    | In order to  create cron job, you need to specify the path to the php executable (default is given for ubuntu)
    | you also need to specify path to artisan, which lies in the root of this repo.
    |
    */

	'php_path' => "/usr/bin/php",
	'artisan_path'=> base_path('artisan'),

	/*
    |--------------------------------------------------------------------------
    | Your AWS Credentials
    |--------------------------------------------------------------------------
    |
    | In order to communicate with an AWS service, you must provide your AWS
    | credentials including your AWS Access Key ID and your AWS Secret Key. You
    | can obtain these keys by logging into your AWS account and visiting
    | https://console.aws.amazon.com/iam/home?#security_credential.
    |
    | Its recommended to create an AWS user with appropriate priveleges and use its access key rather than using a root key.
    |
    */

	'aws_key'      => env('AWS_KEY_ID'),
	'aws_secret'   => env('AWS_KEY_SECRET'),

	/*
    |--------------------------------------------------------------------------
    | AWS Region
    |--------------------------------------------------------------------------
    |
    | Many AWS services are available in multiple regions. You should specify
    | the AWS region you would like to use.
    |
    | These are the regions: us-east-1, us-west-1, us-west-2, us-gov-west-1
    | eu-west-1, sa-east-1, ap-northeast-1, ap-southeast-1, ap-southeast-2
    |
    */

	'aws_region'=>"us-east-1",


	/*
	|--------------------------------------------------------------------------
    | Mail Configuration
    |--------------------------------------------------------------------------
	| This application sends notification mail for all leases created/terminated.
	| Please provide an emailid for sending the mails.
	| In developemnt environment change mail_pretend to true to skip actual sedning of mails & just log in it in laravel log
	| You can also set the global form address & name for the notfication mail
	|
	*/

	'notification_emailid' => "root@localhost",
	'mail_pretend'         => false,
	'mail_from_email'      => "root@aws-sec-manager",
	'mail_from_name'       => "Concierge [AWS Lease Manager]",
);
