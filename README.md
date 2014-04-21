AWS Security Group Manager
=========

* Created as replacement of Dome9 & Cloudpassage for managing & secure access control to EC2 Instances on AWS. 
* Allows you to create secure lease for accessing any port on your AWS security groups. Also, supports guest leasing to grant access to guests to specific services for specific time.
* Uses two factor authentication using duosecurity ( http://duosecurity.com ) for better access control.
* Maintains a log of all the leases acquired/terminated and sends notification mails for them.


Installation
--------------
* Basic installation
```sh
git clone [git-repo-url]
chmod -R o+wx app/storage/
```
* Copy the `app/config/custom_config.sample.php` to `app/config/custom_config.php` and modify appropriately. 
* Copy `app/config/database.sample.php` to `app/config/database.php` and update the password. Also, modify `database/seeds/DatabaseSeeder.php` to add inital set of users.
* Migrate the db with `php artisan migrate` &  seed with `php artisan db:seed`
* Create a symlink to public director in `/var/www` folder or create an apache vhost for the site & open in your browser to test it.
* Create the crontab for correcting leasing by running `php artisan custom:croncreator`. 
* You can check if the crontab can run smoothly by running `php artisan custom:leasemanger` to check it returns no errors.