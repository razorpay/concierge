#!/bin/sh
ARGS=$2

conciergeServer() {
    /concierge/concierge $ARGS &
}

ouath2Server() {
    /go/bin/oauth2_proxy --config=/concierge/oauth2_proxy.cfg
}

if [ "$1" == "start" ]
then
    conciergeServer
    ouath2Server 
else
    echo "Exiting!"
fi