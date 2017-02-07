<?php

/*
 | ----------------------------------------------------------------------------------
 | Detect The Application Environment
 | ----------------------------------------------------------------------------------
 |
 */
use Dotenv\Dotenv;
use Dotenv\Exception\InvalidPathException;

$envDir = __DIR__.'/../environment';

$app->useEnvironmentPath($envDir);

//
// By default we assume environment is prod.
// During testing, laravel sets APP_ENV to 'testing'
// Otherwise, we get the environement from the file
// environment/env.php
//

$env = 'production';

if (env('APP_ENV') === 'testing')
{
    $env = 'testing';
}
else if (file_exists($file = __DIR__ . '/../environment/env.php'))
{
    $env = require $file;
}

putenv("APP_ENV=$env");

$file = $app->environmentFile();

$cascadingEnvFile = '.env.' . $env;

//
// Environment variable files are loaded in the order
// * Vault env file
// * Cascaded environment based env file
// * Default env file
//
// Note that of the above 3, first two are committed in git
// while last one comes into the folder when baking amis via brahma
//

if (! function_exists('read_env_file'))
{
    function read_env_file($envDir, $fileName)
    {
        $file = $envDir . '/' . $fileName;

        if (file_exists($file) === false)
        {
            return;
        }

        $dotenv = new Dotenv($envDir, $fileName);

        $dotenv->load();
    }
}

read_env_file($envDir, '.env.vault');
read_env_file($envDir, $cascadingEnvFile);
read_env_file($envDir, '.env.defaults');
