<?php

return [
    /*
    |--------------------------------------------------------------------------
    | Paths
    |--------------------------------------------------------------------------
    |
    | In order to  create cron job, you need to specify the path to the php executable (default is given for ubuntu)
    | you also need to specify path to artisan, which lies in the root of this repo.
    |
    */

    'php_path'    => '/usr/bin/php',
    'artisan_path'=> base_path('artisan'),

    /*
    |--------------------------------------------------------------------------
    | Mail Configuration
    |--------------------------------------------------------------------------
    | This application sends notification mail for all leases created/terminated.
    | Please provide an emailid for sending the mails.
    | In developemnt environment change mail_pretend to true to skip actual sedning of mails
    | & just log in it in laravel log
    | You can also set the global form address & name for the notfication mail
    |
    */

    'cron_password' => env('CRON_PASSWORD'),
    'google_domain' => env('COMPANY_DOMAIN'),
];
