<?php
return array(

    /*
    |--------------------------------------------------------------------------
    | oAuth Config
    |--------------------------------------------------------------------------
    */

    /**
     * Storage
     */
    'storage' => 'Session',

    /**
     * Consumers
     */
    'consumers' => array(

        /**
         * Google
         */
        'Google' => array(
            'client_id'     => $_ENV['google_client_id'],
            'client_secret' => $_ENV['google_client_secret'],
            'scope'         => array('userinfo_email', 'userinfo_profile'),
        ),

    ),

    'userinfo_url' => 'https://www.googleapis.com/oauth2/v1/userinfo'

);
