#!/bin/sh

CONFIG_FILE=${QDROUTERD_HOME}/etc/qdrouterd.conf

${QDROUTERD_HOME}/bin/configure.sh ${QDROUTERD_HOME} $CONFIG_FILE

if [ -f $CONFIG_FILE ]; then
    ARGS="-c $CONFIG_FILE"
fi

exec qdrouterd $ARGS
