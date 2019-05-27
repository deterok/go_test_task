#!/bin/sh

case "$OSTYPE" in
    darwin*) DOCKER_PATH=$(dirname `readlink $0`);;
    *) DOCKER_PATH=$(dirname `realpath $0`);;
esac

ENV_FILE="$DOCKER_PATH/docker-compose.env"

if [ -f $ENV_FILE ];
    then
        source $ENV_FILE
fi

if [ -z ${DC_ARGS+x} ];
    then
        docker-compose $@
    else
        docker-compose $DC_ARGS $@
fi
