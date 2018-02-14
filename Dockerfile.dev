FROM razorpay/pithos:rzp-php7.1-nginx

COPY . /app/

COPY ./dockerconf/entrypoint.sh /entrypoint.sh

WORKDIR /app

RUN apk update && \
    apk add --no-cache \
    php7-tokenizer \
    php7-xmlwriter \
    php7-simplexml \
    && composer install --no-interaction \
    && composer clear-cache

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
