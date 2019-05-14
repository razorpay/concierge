FROM razorpay/pithos:rzp-php7.1-nginx

LABEL maintainer="Nemo <n@rzp.io>"

ARG GIT_COMMIT_HASH
ENV GIT_COMMIT_HASH=${GIT_COMMIT_HASH}

COPY . /app/

COPY ./dockerconf/entrypoint.sh /entrypoint.sh
COPY ./dockerconf/concierge.nginx.conf /etc/nginx/conf.d/default.conf

RUN /app/dockerconf/build.sh

WORKDIR /app

RUN apk update && \
    apk add --no-cache \
    php7-tokenizer \
    php7-xmlwriter \
    php7-simplexml \
    && composer install --no-interaction --no-dev --optimize-autoloader\
    && composer clear-cache

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
