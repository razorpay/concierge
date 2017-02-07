<?php

return [

    'defaults' => [
        'guard' => 'user',
        'passwords' => 'user',
    ],

    //Authenticating guards
    'guards' => [
        'user' =>[
            'driver' => 'session',
            'provider' => 'users',
        ]
    ],

    //User Providers

    'providers' => [
        'users' => [
            'driver' => 'eloquent',
            'model' => App\Models\User::class,
        ]
    ]

];
