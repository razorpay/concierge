AWS Security Group Manager
=========

Created as replacement of Dome9 & Cloudpassage for managing & secure access control to EC2 Instances

Installation
--------------
* Basic installation
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
* Copy `app/LaravelDuo/LaravelDuo.sample.php` to `app/LaravelDuo/LaravelDuo.php` and update it with necessary values. 
* Copy `app/config/database.sample.php` to `app/config/database.php` and update the password
* Migrate the db with `php artisan migrate` &  seed with `php artisan db:seed`
* Give permissions on storage folder `chmod -R o+wx app/storage/`
* Copy the `app/custom_config.sample.php` to `app/custom_config.php` and modify appropriately.
* Create the crontab for correcting leasing by running `php artisan custom:croncreator`. You can check if the crontab can run smoothly by running `php artisan custom:leasemanger` to check it returns no errors.