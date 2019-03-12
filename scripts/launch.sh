#!/bin/sh

export HOSTNAME_IP_ADDRESS=$(hostname -i)

CONFIG_FILE=/tmp/qdrouterd.conf

${QDROUTERD_HOME}/bin/configure.sh ${QDROUTERD_HOME} $CONFIG_FILE

if [ -f $CONFIG_FILE ]; then
    ARGS="-c $CONFIG_FILE"
fi

exec qdrouterd $ARGS
