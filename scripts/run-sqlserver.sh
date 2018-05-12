#!/bin/bash

set -xeu

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

docker run -t -i \
--workdir /tmp/build \
-e DB_SCHEME=sqlserver \
-e DB_NAME=uaa \
-e DB_USERNAME=root \
-e DB_PASSWORD=changemeCHANGEME1234! \
-e DB_HOSTNAME=127.0.0.1 \
-e DB_PORT=1433 \
-e UAA_DIR=/tmp/build/uaa \
-v ~/workspace/uaa:/tmp/build/uaa \
-v $SCRIPT_DIR/..:/tmp/build/uaa-key-rotator \
cfidentity/uaa-key-rotator-sqlserver \
/tmp/build/uaa-key-rotator/ci/tasks/run-unit-tests-sqlserver/task.sh