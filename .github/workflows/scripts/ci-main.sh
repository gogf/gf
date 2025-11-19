#!/usr/bin/env bash

coverage=$1

# find all path that contains go.mod.
for file in `find . -name go.mod`; do
    dirpath=$(dirname $file)
    echo $dirpath

    # ignore mssql tests as its docker service failed
    # TODO remove this ignoring codes after the mssql docker service OK
    if [ "mssql" = $(basename $dirpath) ]; then
        continue 1
    fi

    # package kubecm was moved to sub ci procedure.
    if [ "kubecm" = $(basename $dirpath) ]; then
        continue 1
    fi

    # examples directory was moved to sub ci procedure.
    if [[ $dirpath =~ "/examples/" ]]; then
        continue 1
    fi

    # Check if it's a contrib directory
    if [[ $dirpath =~ "/contrib/" ]]; then
        # Check if go version meets the requirement
        if ! go version | grep -qE "go${LATEST_GO_VERSION}"; then
            echo "ignore path $dirpath as go version is not ${LATEST_GO_VERSION}: $(go version)"
            continue 1
        fi
    fi

    if [[ $file =~ "/testdata/" ]]; then
        echo "ignore testdata path $file"
        continue 1
    fi

    if [[ $dirpath = "." ]]; then
        # No space left on device error sometimes occurs in CI pipelines, so clean the cache before tests.
        go clean -cache
        # docker stop $(docker ps -aq)
        # docker rm $(docker ps -aq)
        # docker rmi -f $(docker images -aq)
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
