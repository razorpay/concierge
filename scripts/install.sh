#!/bin/bash
# Deployment Script
BASEDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )/../" && pwd )

# Take the app down
cd /home/ubuntu/concierge/ && php artisan down

# Install new version
cd $BASEDIR && rsync -avz --force --delete --progress --exclude-from=./.rsyncignore ./ /home/ubuntu/concierge/

# Fix permissions
cd /home/ubuntu/concierge/ && sudo chmod 775 -R app/storage

# DB Migrate 
cd /home/ubuntu/concierge/ && php artisan migrate

# Take the app up
cd /home/ubuntu/concierge/ && php artisan up