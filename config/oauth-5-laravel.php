<?php

return [

	/*
	|--------------------------------------------------------------------------
	| oAuth Config
	|--------------------------------------------------------------------------
	*/

	/**
	 * Consumers
	 */
	'consumers' => [

        /**
         * Google
         */
        'Google' => [
            'client_id'     => env('google_client_id'),
            'client_secret' => env('google_client_secret'),
            'scope'         => ['userinfo_email', 'userinfo_profile'],
        ],

	],

    'userinfo_url' => 'https://www.googleapis.com/oauth2/v1/userinfo'

];
