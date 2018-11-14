FROM alpine:3.7
COPY composer.json composer.lock /app/

LABEL maintainer="Nemo <n@rzp.io>"
#TODO: Not needed
ENV APACHE_DOCUMENT_ROOT /app/public

WORKDIR /app

RUN apk add --no-cache \
    php7 php7-curl php7-openssl php7-pdo php7-mbstring php7-tokenizer php7-xml \
    php7-mysqlnd php7-pdo_mysql php7-pdo_sqlite php7-fpm php7-session php7-phar \
    php7-simplexml php7-ctype php7-json nginx \
    && mkdir -p /run/nginx/ && \
    chown 0775 /run/nginx/

RUN wget https://raw.githubusercontent.com/composer/getcomposer.org/1b137f8bf6db3e79a38a5bc45324414a6b1f9df2/web/installer -O - -q | php -- --quiet

RUN composer install --no-dev --no-interaction --no-autoloader --no-scripts && \
    rm -rf /root/.composer && \
    composer clear-cache

COPY dockerconf/nginx.conf /etc/nginx/
COPY dockerconf/default.conf /etc/nginx/conf.d/
COPY dockerconf/php-fpm.conf /etc/php7/
COPY dockerconf/www.conf /etc/php7/php-fpm.d

COPY . /app/

RUN composer dump-autoload -o && php artisan optimize && chmod -R o+wx storage

EXPOSE 80

ENTRYPOINT /app/dockerconf/entrypoint.sh
