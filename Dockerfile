FROM razorpay/containers:rzp-php7.1-nginx

ARG GIT_COMMIT_HASH
ENV GIT_COMMIT_HASH=${GIT_COMMIT_HASH}

RUN mkdir /app/

COPY . /app/

COPY ./dockerconf/entrypoint.sh /entrypoint.sh

RUN /app/dockerconf/build.sh

WORKDIR /app

RUN composer install --no-interaction

EXPOSE 80

ENTRYPOINT ["/usr/local/bin/dumb-init", "--"]
CMD ["/entrypoint.sh"]
