#!/bin/bash
# Deployment Script
echo "Setting BASEDIR"
BASEDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )/../" && pwd )

# Take the app down
echo "Take the app down"
cd /home/ubuntu/concierge/ && php artisan down

# Install new version
echo  "Install new version"
cd $BASEDIR && rsync -avz --force --delete --progress --exclude-from=./.rsyncignore ./ /home/ubuntu/concierge/

# DB Migrate
echo  "DB Migrate"
cd /home/ubuntu/concierge/ && php artisan migrate

# Take the app up
echo  "Take the app up"
cd /home/ubuntu/concierge/ && php artisan up
