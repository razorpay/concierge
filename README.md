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
* Create a new user in AWS IAM. Note its access key & secret for next step. And attach following policy to it
```
	{
	  "Statement": [
	    {
	      "Action": [
	        "ec2:AuthorizeSecurityGroupEgress",
	        "ec2:AuthorizeSecurityGroupIngress",
	        "ec2:CreateSecurityGroup",
	        "ec2:DeleteSecurityGroup",
	        "ec2:DescribeInstances",
	        "ec2:DescribeSecurityGroups",
	        "ec2:RevokeSecurityGroupEgress",
	        "ec2:RevokeSecurityGroupIngress"
	      ],
	      "Effect": "Allow",
	      "Resource": "*"
	    }
	  ]
	}
```
* Copy the `app/config/custom_config.sample.php` to `app/config/custom_config.php` and modify appropriately. 
* Copy `app/config/database.sample.php` to `app/config/database.php` and update the password. 
* Migrate the db with `php artisan migrate` &  seed with `php artisan db:seed`
* Create a symlink to public director in `/var/www` folder or create an apache vhost for the site & open in your browser to test it.
* Create the crontab for managing expired leases by running `php artisan custom:croncreator`. 
* You can check if the crontab can run smoothly by running `php artisan custom:leasemanger` to check it returns no errors.