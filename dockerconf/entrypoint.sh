#!/bin/bash
set -euo pipefail

# wait for db to be provisioned
sleep 30

echo "Creating Log Files"
mkdir -p /var/log/nginx

ALOHOMORA_BIN=$(which alohomora)

# fix permissions
chown -R nginx.nginx /app

# TODO: merge concierge.nginx.conf & concierge.nginx.dev.conf in a single j2 template
# once `alohomora` has simple jinja2 render capabilities
# without credstash lookups

if [[ "${APP_ENV}" == "dev" ]]; then
  cp environment/.env.docker environment/.env.testing
  cp environment/env.sample.php environment/env.php
  sed -i 's/dev/testing/g' ./environment/env.php
  cp dockerconf/concierge.dev.nginx.conf /etc/nginx/conf.d/concierge.conf
else
  # casting alohomora to unlock the secrets
  $ALOHOMORA_BIN cast --region ap-south-1 --env $APP_ENV --app concierge "environment/.env.j2"
  $ALOHOMORA_BIN cast --region ap-south-1 --env $APP_ENV --app concierge "environment/env.php.j2"
  $ALOHOMORA_BIN cast --region ap-south-1 --env $APP_ENV --app concierge "dockerconf/concierge.nginx.conf.j2"
  cp dockerconf/concierge.nginx.conf /etc/nginx/conf.d/concierge.conf
fi

echo "create /app/storage/framework/views/"
mkdir storage/framework/views

# Fix permissions
echo  "$(date) Fix Storage permissions"
cd /app/ && chmod -R o+wx storage/

mkdir -p /run/
chown 0755 /run/

# for enabling reading off env vars from os
sed -i -- 's/;clear_env = no/clear_env = no/g' /etc/php7/php-fpm.d/www.conf

echo "GIT_COMMIT_HASH=${GIT_COMMIT_HASH}" >> /app/.env

/usr/sbin/php-fpm7
/usr/sbin/nginx -g 'daemon off;'
