FROM razorpay/containers:base-nginx-php7

ARG GIT_COMMIT_HASH
ENV GIT_COMMIT_HASH=${GIT_COMMIT_HASH}

COPY . /app/

COPY ./dockerconf/entrypoint.sh /entrypoint.sh

RUN /app/dockerconf/build.sh

WORKDIR /app

RUN composer install --no-interaction

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
