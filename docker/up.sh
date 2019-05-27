#!/bin/sh

DEPLOY_PATH=$(dirname `realpath $0`)

$DEPLOY_PATH/dc.sh up $@
