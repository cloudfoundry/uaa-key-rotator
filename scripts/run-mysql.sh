#!/bin/bash

set -xeu

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

docker run -t -i \
--workdir /tmp/build \
-e DB_SCHEME=mysql \
-e DB_NAME=uaa \
-e DB_USERNAME=root \
-e DB_PASSWORD=changeme \
-e DB_HOSTNAME=localhost \
-e DB_PORT=3306 \
-v ~/workspace/uaa:/tmp/build/uaa \
-v $SCRIPT_DIR/..:/tmp/build/uaa-key-rotator \
cfidentity/uaa-key-rotator-mysql \
/tmp/build/uaa-key-rotator/ci/tasks/run-unit-tests-mysql/task.sh