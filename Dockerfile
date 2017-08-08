FROM razorpay/containers:base-nginx-php7

COPY . /app/

RUN chown -R nginx.nginx /app

COPY ./dockerconf/entrypoint.sh /entrypoint.sh

WORKDIR /app

RUN chmod +x /entrypoint.sh && \
    /usr/local/bin/composer install --no-interaction

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
