#!/bin/bash
set -euo pipefail

# wait for db to be provisioned
sleep 30

echo "Creating Log Files"
mkdir -p /var/log/nginx

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

  # Seed the database
  # We don't do this for prod, because DB_* is expected to point to RDS instead
  cat database/seed.sql | mysql -u "$DB_USERNAME" -h "$DB_HOST" -P 3306 --password="$DB_PASSWORD" "$DB_DATABASE"

  # for enabling reading off env vars from os since Docker passes these vars
  # to nginx, which can then pass them to PHP/FPM
  sed -i -- 's/;clear_env = no/clear_env = no/g' /etc/php7/php-fpm.d/www.conf

else
  ALOHOMORA_BIN=$(which alohomora)
  # casting alohomora to unlock the secrets
  echo "<?php return 'production';" > "environment/env.php"
  sed -i "s|NGINX_HOST|$HOSTNAME|g" dockerconf/concierge.nginx.conf.j2
  $ALOHOMORA_BIN cast --region ap-south-1 --env $APP_ENV --app concierge "dockerconf/concierge.nginx.conf.j2"
  $ALOHOMORA_BIN cast --region ap-south-1 --env $APP_ENV --app concierge "dockerconf/fastcgi.conf.j2"
  cp dockerconf/concierge.nginx.conf /etc/nginx/conf.d/concierge.conf
  cp dockerconf/fastcgi.conf /etc/nginx/
fi

# Fix permissions
echo  "$(date) Fix Storage permissions"
cd /app/ && chmod -R o+wx storage/

mkdir -p /run/
chown 0755 /run/

/usr/sbin/php-fpm7
/usr/sbin/nginx -g 'daemon off;'
