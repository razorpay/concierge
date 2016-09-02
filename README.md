Concierge - AWS EC2 Secure Lease Manager
=========
Web Management & Secure Access Control of AWS Security Groups. Created as replacement of Dome9 & Cloudpassage for managing & secure access control to EC2 Instances on AWS. Allows you to create IP leases for incoming connections to AWS instances. Also, allows creating invites to allow guest access by URL/Email Invites. Leases expire after specific time or can be manually terminated.

Features
---------
* Allows you to create secure lease for accessing any port on your AWS security groups.
* Supports guest leasing to grant access to guests to specific services for specific time via URL or Email invite.
* Support for deploy Invites, that can be used in auto-deployment scripts from wercker, shippable etc. Deploy invites can be used repeatedly unlike URL invites which are single use.
* Uses two factor authentication using duosecurity ( http://duosecurity.com ) for better access control.
* Maintains a record of all the leases acquired/terminated in the database and sends notification mails for them.
* Maintains a log file of the leases.



Installation
--------------
* Basic installation
```sh
git clone [git-repo-url]
chmod -R o+wx app/storage/
```
* Create a new user in AWS IAM. Note its access key & secret for next step. And attach following policy to it
```
	{
	  "Statement": [
	    {
	      "Action": [
	        "ec2:AuthorizeSecurityGroupIngress",
	        "ec2:DescribeSecurityGroups",
	        "ec2:RevokeSecurityGroupIngress"
	      ],
	      "Effect": "Allow",
	      "Resource": "*"
	    }
	  ]
	}
```
* Create `bootstrap/env.php` with your environment (`local` for dev).
* Copy `env.sample.php` to `.env.[production|local].php` and add your google auth client id/secret.
* Copy the `app/config/custom_config.sample.php` to `app/config/custom_config.php` and modify appropriately.
* Copy `app/config/database.sample.php` to `app/config/database.php` and update the password.
* Install laravel with `php composer.phar install`. Migrate the db with `php artisan migrate` &  seed with `php artisan db:seed`.
* Create a symlink to "public" directory in `/var/www` folder or create an apache vhost for the site & open in your browser to test it.
* Create the crontab for managing expired leases by running `php artisan custom:croncreator`.
* You can check if the crontab can run smoothly by running `php artisan custom:leasemanger` to check it returns no errors.

Screenshots
--------------
### Home Page ![Home](/screenshots/home.png?raw=true "Home Page")

### Manage Security Group ![Manage Groups](/screenshots/manage-group.png?raw=true "Manage Group")

### New URL Invite ![URL Invite](/screenshots/url-invite.png?raw=true "New URL Invite")

### Manage Users ![Manage Users](/screenshots/manage-users.png?raw=true "Manage Users")

License
--------------
This application is open-sourced software licensed under the MIT license
