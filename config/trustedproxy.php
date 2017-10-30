<?php
return [
    /*
     * Set trusted proxy IP addresses.
     *
     * Both IPv4 and IPv6 addresses are
     * supported, along with CIDR notation.
     *
     * The "*" character is syntactic sugar
     * within TrustedProxy to trust any proxy;
     * a requirement when you cannot know the address
     * of your proxy (e.g. if using Rackspace balancers).
     */
    'proxies' => [
        '10.0.0.0/8',
    ],
    /*
     * Default Header Names
     *
     * Change these if the proxy does
     * not send the default header names.
     *
     * Note that headers such as X-Forwarded-For
     * are transformed to HTTP_X_FORWARDED_FOR format.
     *
     * The following are Symfony defaults, found in
     * \Symfony\Component\HttpFoundation\Request::$trustedHeaders
     *
     * You may optionally set headers to 'null' here if you'd like
     * for them to be considered untrusted instead. Ex:
     *
     * Illuminate\Http\Request::HEADER_CLIENT_HOST  => null,
     *
     * WARNING: If you're using AWS Elastic Load Balancing or Heroku,
     * the FORWARDED and X_FORWARDED_HOST headers should be set to null
     * as they are currently unsupported there.
     */
    'headers' => [
        Symfony\Component\HttpFoundation\Request::HEADER_CLIENT_IP    => 'X_FORWARDED_FOR',
        Symfony\Component\HttpFoundation\Request::HEADER_CLIENT_PROTO => 'X_FORWARDED_PROTO',
        Symfony\Component\HttpFoundation\Request::HEADER_CLIENT_PORT  => 'X_FORWARDED_PORT',
        Symfony\Component\HttpFoundation\Request::HEADER_FORWARDED    => null,
        Symfony\Component\HttpFoundation\Request::HEADER_CLIENT_HOST  => null,
    ]
];
