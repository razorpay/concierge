FROM razorpay/onggi:php-7.3-nginx

ARG APP_ENV
ENV APP_ENV=${APP_ENV}

WORKDIR /app

COPY composer.* ./
COPY . .

RUN composer install --no-interaction --no-dev --optimize-autoloader \
    && composer clear-cache

EXPOSE 80

RUN chmod +x ./dockerconf/entrypoint.sh

ENTRYPOINT [ "./dockerconf/entrypoint.sh" ]
