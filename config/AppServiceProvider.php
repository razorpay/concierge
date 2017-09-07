<?php

namespace App\Providers;

use SocialOAuth;
use Illuminate\Support\ServiceProvider;

class AppServiceProvider extends ServiceProvider
{
    protected $defer = false;
    /**
     * Bootstrap any application services.
     *
     * @return void
     */
    public function boot()
    {
        $this->configureOAuth();
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
