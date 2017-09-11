<?php

namespace App\Providers;

use SocialOAuth;
use Illuminate\Support\ServiceProvider;
use Illuminate\Http\Request;

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
        $this->configureTrustedProxies($request);
    }

    protected function configureTrustedProxies($request)
    {
        foreach (config('trustedproxy.headers') as $headerKey => $headerName)
        {
            $request->setTrustedHeaderName($headerKey, $headerName);
        }

        $request->setTrustedProxies(config('trustedproxy.proxies'));
    }

    /**
     * The default is set to StreamClient
     * which is horrible in performance
     */
    protected function configureOAuth()
    {
        SocialOAuth::setHttpClient('CurlClient');
    }

    public function register()
    {
        ;
    }
}
