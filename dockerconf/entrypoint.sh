#!/bin/sh
set -euo pipefail

cp .env.example .env
php artisan key:generate

php artisan optimize:clear
php artisan route:clear
php artisan config:clear
php artisan cache:clear

chmod -R 777 storage/

# Waiting for DB
chmod +x ./dockerconf/wait_for_db.sh
./dockerconf/wait_for_db.sh

if $CRON_ENABLE
then
    echo "Enabling Cron to remove expired/outdated leases";
    php artisan concierge:cronjob
fi

cp ./dockerconf/concierge.nginx.conf /etc/nginx/conf.d/concierge.conf
chown -R nginx:nginx /app/storage/logs

/usr/sbin/php-fpm7
/usr/sbin/nginx -g 'daemon off;'
