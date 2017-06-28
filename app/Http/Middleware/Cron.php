<?php

namespace App\Http\Middleware;

use Auth;
use Closure;
use Illuminate\Http\Response;

class Cron
{
    public function __construct()
    {
    }

    /**
     * Handle an incoming request.
     *
     * @param \Illuminate\Http\Request $request
     * @param \Closure                 $next
     *
     * @return mixed
     */
    public function handle($request, Closure $next)
    {
        $username = 'cron';
        $password = $request->getPassword();

        if ($request->getUser() === $username and password_verify($password, config('concierge.cron_password')))
        {
            return $next($request);
        }

        return new Response('Invalid credentials.', 401, ['WWW-Authenticate' => 'Basic']);
    }
}
