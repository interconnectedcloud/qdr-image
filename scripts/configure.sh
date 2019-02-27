#!/bin/bash

HOME_DIR=$1
OUTFILE=$2
INFILE=$HOME_DIR/etc/qdrouterd.conf.in

function substVars() {
    DOLLAR='$' envsubst < $1 > $2
}

function printConfig() {
    echo "---------------------------------------" && cat $OUTFILE && echo "---------------------------------------"
}

if [[ $QDROUTERD_CONF =~ .*\{.*\}.* ]]; then
    # env var contains inline config
    echo "$QDROUTERD_CONF" > $INFILE
elif [[ -n $QDROUTERD_CONF ]]; then
    # treat as path(s)
    IFS=':,' read -r -a array <<< "$QDROUTERD_CONF"
    > $INFILE
    for i in "${array[@]}"; do
        if [[ -d $i ]]; then
            # if directory, concatenate to output all .conf files
            # within it
            for f in $i/*.conf; do
                cat "$f" >> $INFILE
            done
        elif [[ -f $i ]]; then
            # if file concatenate that to the output
            cat "$i" >> $INFILE
        else
            echo "No such file or directory: $i"
        fi
    done
fi

if [ -f $INFILE ]; then
    substVars $INFILE $OUTFILE
fi

if [ -n "$QDROUTERD_AUTO_MESH_DISCOVERY" ]; then
    python $HOME_DIR/bin/auto_mesh.py $OUTFILE || printConfig
fi

if [ -n "$QDROUTERD_AUTO_CREATE_SASLDB_SOURCE" ]; then
    $HOME_DIR/bin/create_sasldb.sh ${QDROUTERD_AUTO_CREATE_SASLDB_PATH:-$HOME_DIR/etc/qdrouterd.sasldb} $QDROUTERD_AUTO_CREATE_SASLDB_SOURCE "${APPLICATION_NAME:-qdrouterd}"
fi
