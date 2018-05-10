#!/bin/bash

export DB_SCHEME=sqlserver
export DB_NAME=uaa
export DB_USERNAME=root
export DB_PASSWORD=changemeCHANGEME1234!
export DB_HOSTNAME=127.0.0.1
export DB_PORT=1433
export UAA_LOCATION=~/workspace/uaa


ginkgo -r .