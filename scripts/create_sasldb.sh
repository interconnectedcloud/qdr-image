#!/bin/bash

SASLDB=$1
USER_DIR=$2
DOMAIN=$3

rm -rf $SASLDB
for user in $USER_DIR/*; do
    echo "cat $user | saslpasswd2 -c -p -u $DOMAIN $(basename $user) -f $SASLDB"
    cat $user | saslpasswd2 -c -p -u $DOMAIN $(basename $user) -f $SASLDB
done

