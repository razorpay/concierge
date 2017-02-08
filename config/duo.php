<?php

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

return [
    'akey' => env('DUO_AKEY'),
    'ikey' => env('DUO_IKEY'),
    'skey' => env('DUO_SKEY'),
    'host' => env('DUO_HOST'),
]
