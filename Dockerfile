FROM golang:1.12-alpine3.9 as concierge
RUN apk add git
WORKDIR /concierge
COPY go.mod go.sum ./
RUN go mod download
RUN export GO111MODULE=on && go get github.com/pusher/oauth2_proxy && export GO111MODULE=auto
COPY . .
RUN go build -o concierge main.go 
EXPOSE 8990 3306 4180
ENTRYPOINT ["docker/entrypoint.sh"]

# FROM razorpay/onggi:base-3.7
# WORKDIR /app
# RUN wget -O /tmp/oauth2_proxy.tar.gz https://github.com/bitly/oauth2_proxy/releases/download/v2.2/oauth2_proxy-2.2.0.linux-amd64.go1.8.1.tar.gz && \
#     cd /tmp/ && \
#     tar -xzf /tmp/oauth2_proxy.tar.gz oauth2_proxy-2.2.0.linux-amd64.go1.8.1/oauth2_proxy --strip-components=1 && \
#     mv /tmp/oauth2_proxy /usr/local/bin/oauth2_proxy && \
#     chmod +x /usr/local/bin/oauth2_proxy && \
#     rm /tmp/oauth2_proxy.tar.gz
# ENV SRC_DIR=/go/src/github.com/razorpay/credstash-ui
# COPY --from=credstash-ui $SRC_DIR/docker docker
# COPY --from=credstash-ui $SRC_DIR/dist dist
# COPY --from=credstash-ui $SRC_DIR/credstash-ui credstash-ui
# # Disable group authentication
# COPY --from=lomkju/oauth2-proxy /oauth2_proxy /usr/local/bin/oauth2_proxy
# ENTRYPOINT ["docker/entrypoint.sh"]