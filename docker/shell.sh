#!/bin/sh

DEPLOY_PATH=$(dirname `realpath $0`)
CONTAINER_NAME=$1

if [ -z $2 ];
    then
        CMD="/bin/sh"
    else
        CMD=$2
fi

$DEPLOY_PATH/dc.sh exec $CONTAINER_NAME $CMD
