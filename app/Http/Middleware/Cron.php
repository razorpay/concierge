<?php

namespace App\Http\Middleware;

use Closure;

class Cron
{
    /**
     * Handle an incoming request.
     *
     * @param  \Illuminate\Http\Request  $request
     * @param  \Closure  $next
     * @return mixed
     */
    public function handle($request, Closure $next)
    {
        $username = 'cron';
        $password = $request->getPassword();

        if (!is_null($password) and hash_equals($request->getUser(), $username) and
            hash_equals($password, config('concierge.cron_password'))) {
            return $next($request);
        }

        return response()->make('Invalid credentials.', 401, ['WWW-Authenticate' => 'Basic']);
    }
}
