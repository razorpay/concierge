#!/bin/sh
php artisan migrate --force
php-fpm7
nginx -g 'daemon off;'
