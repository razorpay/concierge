#!/bin/bash

# copy nginx config
cp dockerconf/concierge.docker.conf /etc/nginx/conf.d/concierge.conf

cp environment/env.sample.php environment/env.php

echo "create /app/storage/framework/views/"
mkdir storage/framework/views

# Fix permissions
echo  "$(date) Fix Storage permissions"
cd /app/ && chmod -R o+wx storage/

mkdir -p /run/
chown 0755 /run/

sed -i -- 's/;clear_env = no/clear_env = no/g' /etc/php7/php-fpm.d/www.conf

echo "GIT_COMMIT_HASH=${GIT_COMMIT_HASH}" >> /app/.env

/usr/sbin/php-fpm7
/usr/sbin/nginx -g 'daemon off;'
