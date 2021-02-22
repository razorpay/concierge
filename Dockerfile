FROM quay.io/oauth2-proxy/oauth2-proxy:v6.0.0 as oauth2_proxy

FROM golang:1.14.0-alpine3.11 as concierge
# hadolint ignore=DL3018
RUN apk update && \
    apk add --no-cache git && \
    rm /var/cache/apk/*

WORKDIR /concierge
COPY go.mod go.sum ./
RUN GO111MODULE=on go mod download
COPY . .
RUN go build -o concierge main.go 


FROM razorpay/onggi:base-3.7
WORKDIR /app

COPY --from=concierge /concierge/concierge concierge
COPY --from=concierge /concierge/assets assets
COPY --from=concierge /concierge/templates templates
COPY --from=concierge /concierge/docker docker
COPY --from=concierge /concierge/oauth2_proxy oauth2_proxy
COPY --from=oauth2_proxy /bin/oauth2-proxy /usr/local/bin/oauth2-proxy

EXPOSE 8990 3306 4180
ENTRYPOINT ["docker/entrypoint.sh"]
