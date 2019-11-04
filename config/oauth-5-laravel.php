<?php

return [

    /*
    |--------------------------------------------------------------------------
    | oAuth Config
    |--------------------------------------------------------------------------
    */

    /**
     * Storage
     */
    'storage' => '\\OAuth\\Common\\Storage\\Session',

    /**
     * Consumers
     */
    'consumers' => [

        /*
         * Google
         */
        'Google' => [
            'client_id'     => env('GOOGLE_CLIENT_ID'),
            'client_secret' => env('GOOGLE_CLIENT_SECRET'),
            'scope'         => ['userinfo_email', 'userinfo_profile'],
        ],

    ],

    'userinfo_url' => 'https://www.googleapis.com/oauth2/v1/userinfo',

];
