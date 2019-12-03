#!/bin/sh
ARGS=$2

conciergeServer() {
    /app/concierge $ARGS &
}

ouath2Server() {
    /usr/local/bin/oauth2_proxy --config=/app/oauth2_proxy/oauth2_proxy.cfg
}

if [ "$1" == "start" ]
then
    conciergeServer
    ouath2Server 
else
    echo "Exiting!"
fi