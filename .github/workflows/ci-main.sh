#!/usr/bin/env bash

# Define the latest Go version requirement
LATEST_GO_VERSION="1.23"

coverage=$1

# find all path that contains go.mod.
for file in `find . -name go.mod`; do
    dirpath=$(dirname $file)
    echo $dirpath

    # Check if it's a contrib directory or example directory
    if [[ $dirpath =~ "/contrib/" ]] || [ "example" = $(basename $dirpath) ]; then
        # Check if go version meets the requirement
        if ! go version | grep -qE "go${LATEST_GO_VERSION}"; then
            echo "ignore path $dirpath as go version is not ${LATEST_GO_VERSION}: $(go version)"
            continue 1
        fi
        # If it's example directory, only build without tests
        if [ "example" = $(basename $dirpath) ]; then
            echo "the example directory only needs to be built, not unit tests and coverage tests."
            cd $dirpath
            go mod tidy
            go build ./...
            cd -
            continue 1
        fi
    fi

    if [[ $file =~ "/testdata/" ]]; then
        echo "ignore testdata path $file"
        continue 1
    fi

    cd $dirpath
    go mod tidy
    go build ./...

    # test with coverage
    if [ "${coverage}" = "coverage" ]; then
      go test ./... -race -coverprofile=coverage.out -covermode=atomic -coverpkg=./...,github.com/gogf/gf/... || exit 1

      if grep -q "/gogf/gf/.*/v2" go.mod; then
        sed -i "s/gogf\/gf\(\/.*\)\/v2/gogf\/gf\/v2\1/g" coverage.out
      fi
    else
      go test ./... -race || exit 1
    fi

    cd -
done
