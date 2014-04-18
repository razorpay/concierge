AWS Security Group Manager
=========

Created as replacement of Dome9 & Cloudpassage for managing & secure access control to EC2 Instances

Installation
--------------
*Basic installation
```sh
git clone [git-repo-url]
chmod -R o+wx storage/
php artisan config:publish aws/aws-sdk-php-laravel
```
* Update your settings in the generated `app/config/packages/aws/aws-sdk-php-laravel/config.php` configuration file.

```
return array(
    'key'         => 'YOUR_AWS_ACCESS_KEY_ID',
    'secret'      => 'YOUR_AWS_SECRET_KEY',
    'region'      => 'us-east-1',
    'config_file' => null,
);
```
* Create a symlink to public director in `/var/www` folder or create an apache vhost for the site & open in your browser.


    