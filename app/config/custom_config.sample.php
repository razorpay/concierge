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
	'php_path'=>"/usr/bin/php",
	'artisan_path'=>"/path/to/aws-secmanager/artisan",


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
	'aws_key'=>"AWS KEY ID",
	'aws_secret'=>"AWS KEY SECRET",

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
	'aws_region'=>"AWS REGION",


	/* 
	|--------------------------------------------------------------------------
    | Duo Two Factor Authentication Configuration
    |--------------------------------------------------------------------------
	| First dign up for duosecurity and create a new Web SDK integration
	| note down the ikey, skey, host from the integration
	| the akey requires a randomly generated string with at least 40 characters, you can use any genertor or string as you wish
	| Make sure to keep the akey secret and only withing one application
	|
	*/
	'duo_akey' => "APPLICATION KEY 45CHARACTER RANDOMLY GENERATED SUPER SECRET",
	'duo_ikey' => "IKEY FROM DUO INTEGRATION",
	'duo_skey' => "SKEY FROM DUO INTEGRATION",
	'duo_host' => "HOST FROM DUO INTEGRATION",

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

	'notification_emailid'=>"root@localhost",
	'mail_pretend' => FALSE,
	'mail_from_email' => "root@aws-sec-manager",
	'mail_from_name' => "AWS Secure Lease Manager",

    /* 
    |--------------------------------------------------------------------------
    | Database Seed
    |--------------------------------------------------------------------------
    | This is used to create default users for the application so that you can log on.
    | You need to only add name, username, password. Also, admin can be TRUE or FALSE for creating site admins or standard users,
    | Admin users can add/delete users & other admins. 
    | No Validation checks are performed on this data so be careful.
    | You need to run "php artisan migrate --seed" for seeding this data.
    |
    */
    'users'=>array( 
            array(
                'name' => 'Admin Name',
                'username' => 'admin',
                'password' => Hash::make('password')
                'admin' => TRUE
            ),
            array(
                'name' => 'Standard User Name',
                'username' => 'user',
                'password' => Hash::make('password'),
                'admin' => FALSE
            ),
            // add more arrays as your requirement
    ),


);