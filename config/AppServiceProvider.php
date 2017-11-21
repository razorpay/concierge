<?php

namespace App\Providers;

use SocialOAuth;
use Illuminate\Http\Request;
use Illuminate\Support\ServiceProvider;
use Illuminate\Support\Facades\Validator;

class AppServiceProvider extends ServiceProvider
{
    protected $defer = false;
    /**
     * Bootstrap any application services.
     *
     * @return void
     */
    public function boot(Request $request)
    {
        $this->configureOAuth();
        $this->configureValidator();
        $this->configureTrustedProxies($request);
    }

    private function configureTrustedProxies($request)
    {
        foreach (config('trustedproxy.headers') as $headerKey => $headerName)
        {
            $request->setTrustedHeaderName($headerKey, $headerName);
        }

        $request->setTrustedProxies(config('trustedproxy.proxies'));
    }

    private function configureValidator()
    {

        Validator::extend('org_email', function ($attribute, $value) {

            $domain = config('concierge.google_domain');

            $email_parts = explode('@', $value);

            $email_domain = end($email_parts);

            return ($email_domain === $domain);
        });

    }

    /**
     * The default is set to StreamClient
     * which is horrible in performance
     */
    private function configureOAuth()
    {
        SocialOAuth::setHttpClient('CurlClient');
    }

    public function register()
    {
        ;
    }
}
