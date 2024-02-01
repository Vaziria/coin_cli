#!/bin/bash

if [[ -z "$COIN_ENTRYPOINT" ]]; then
    COIN_ENTRYPOINT="tonnaged"
    echo "ENTRYPOINT NOT FOUND"
else
    echo "using entrypoint ${COIN_ENTRYPOINT}"
fi


cp /root/coin/src/$COIN_ENTRYPOINT /root/$COIN_ENTRYPOINT 2>/dev/null

if [ -f /root/$COIN_ENTRYPOINT ]; then
    echo "Wallet Exist."
    "/root/${COIN_ENTRYPOINT}"
else
    echo "Wallet NOT EXIST. sleep infinity"
    sleep infinity
fi