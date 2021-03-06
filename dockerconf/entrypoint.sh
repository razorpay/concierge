#!/bin/bash
set -euo pipefail

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

  # Make sure cron runs successfully
  if [[ "${APP_MODE}" == "cron" ]]; then
    echo "Running in cron mode"
    $ALOHOMORA_BIN cast --region ap-south-1 --env $APP_ENV --app concierge "dockerconf/env.sh.j2"
    source "dockerconf/env.sh"
    php artisan custom:leasemanager
    echo "Cron run successfully"
    exit 0
  fi

  $ALOHOMORA_BIN cast --region ap-south-1 --env $APP_ENV --app concierge "dockerconf/fastcgi.conf.j2"
  cp dockerconf/fastcgi.conf /etc/nginx/
  chown -R nginx:nginx /app/storage/logs
fi

/usr/sbin/php-fpm7
/usr/sbin/nginx -g 'daemon off;'
