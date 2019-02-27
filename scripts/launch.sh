#!/bin/sh

export HOSTNAME_IP_ADDRESS=$(hostname -i)

CONFIG_FILE=${QDROUTERD_HOME}/etc/qdrouterd.conf
CONFIG_MAP_FILE=/etc/qpid-dispatch/qdrouterd.conf.template

if [ -f $CONFIG_MAP_FILE ]; then
    DOLLAR='$' envsubst < $CONFIG_MAP_FILE > $CONFIG_FILE
fi

${QDROUTERD_HOME}/bin/configure.sh ${QDROUTERD_HOME} $CONFIG_FILE

if [ -f $CONFIG_FILE ]; then
    ARGS="-c $CONFIG_FILE"
fi

exec qdrouterd $ARGS
