#!/bin/bash

export UAA_LOCATION=$PWD/uaa

source uaa/scripts/start_db_helper.sh
bootDB $DB_SCHEME

set -eux

export GOPATH=$PWD/go
export PATH=$PATH:$GOPATH/bin

mkdir -p go/src/github.com/cloudfoundry/uaa-key-rotator
cp -R uaa-key-rotator/* go/src/github.com/cloudfoundry/uaa-key-rotator
cd go/src/github.com/cloudfoundry/uaa-key-rotator

go get github.com/onsi/ginkgo/ginkgo

reformatted_packages=$(go fmt github.com/cloudfoundry/uaa-key-rotator/...)
if [[ $reformatted_packages = *[![:space:]]* ]]; then
  echo "FAILURE: go fmt reformatted the following packages:"
  echo $reformatted_packages
  exit 1
fi


ginkgo -v -r --race -randomizeAllSpecs -randomizeSuites .