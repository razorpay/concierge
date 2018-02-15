FROM razorpay/pithos:rzp-php7.1-nginx

ARG GIT_COMMIT_HASH
ENV GIT_COMMIT_HASH=${GIT_COMMIT_HASH}

RUN mkdir /app/

COPY . /app/

COPY ./dockerconf/entrypoint.sh /entrypoint.sh

RUN /app/dockerconf/build.sh

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
