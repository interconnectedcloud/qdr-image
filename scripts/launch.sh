#!/bin/sh

export HOSTNAME_IP_ADDRESS=$(hostname -i)

EXT=${QDROUTERD_CONF_TYPE:-conf}
CONFIG_FILE=/tmp/qdrouterd.${EXT}

${QDROUTERD_HOME}/bin/configure.sh ${QDROUTERD_HOME} $CONFIG_FILE

if [ -f $CONFIG_FILE ]; then
    ARGS="-c $CONFIG_FILE"
fi

exec qdrouterd $ARGS
