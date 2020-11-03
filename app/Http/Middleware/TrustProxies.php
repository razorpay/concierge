<?php

namespace App\Http\Middleware;

use Fideloper\Proxy\TrustProxies as Middleware;
use Illuminate\Http\Request;

class TrustProxies extends Middleware
{
    /**
     * The trusted proxies for this application.
     *
     * @var array|string
     *
     * Bug with the working of trusted proxies from v3.3 to v4+ upgrade
     * Open Issue- https://github.com/fideloper/TrustedProxy/issues/115
     * Comment Link- https://github.com/fideloper/TrustedProxy/issues/115#issuecomment-469459016
     */
    protected $proxies = [
        '0.0.0.0/0',
        '2000:0:0:0:0:0:0:0/3'
    ];

    /**
     * The headers that should be used to detect proxies.
     *
     * @var int
     */
    protected $headers = Request::HEADER_X_FORWARDED_ALL;
}
