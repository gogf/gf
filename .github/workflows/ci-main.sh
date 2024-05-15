#!/usr/bin/env bash

coverage=$1

# find all path that contains go.mod.
for file in `find . -name go.mod`; do
    dirpath=$(dirname $file)
    echo $dirpath

    if [[ $file =~ "/testdata/" ]]; then
        echo "ignore testdata path $file"
        continue 1
    fi

    # package kuhecm needs golang >= v1.19
    if [ "kubecm" = $(basename $dirpath) ]; then
        continue 1
        if ! go version|grep -qE "go1.19|go1.[2-9][0-9]"; then
          echo "ignore kubecm as go version: $(go version)"
          continue 1
        fi
    fi

    # package consul needs golang >= v1.19
    if [ "consul" = $(basename $dirpath) ]; then
        continue 1
        if ! go version|grep -qE "go1.19|go1.[2-9][0-9]"; then
          echo "ignore consul as go version: $(go version)"
          continue 1
        fi
    fi

    # package etcd needs golang >= v1.19
    if [ "etcd" = $(basename $dirpath) ]; then
        if ! go version|grep -qE "go1.19|go1.[2-9][0-9]"; then
          echo "ignore etcd as go version: $(go version)"
          continue 1
        fi
    fi

    # package polaris needs golang >= v1.19
    if [ "polaris" = $(basename $dirpath) ]; then
        if ! go version|grep -qE "go1.19|go1.[2-9][0-9]"; then
          echo "ignore polaris as go version: $(go version)"
          continue 1
        fi
    fi

    # package example needs golang >= v1.20
    if [ "example" = $(basename $dirpath) ]; then
        if ! go version|grep -qE "go1.[2-9][0-9]"; then
          echo "ignore example as go version: $(go version)"
          continue 1
        fi
    fi

    # package otlpgrpc needs golang >= v1.20
    if [ "otlpgrpc" = $(basename $dirpath) ]; then
        if ! go version|grep -qE "go1.[2-9][0-9]"; then
          echo "ignore otlpgrpc as go version: $(go version)"
          continue 1
        fi
    fi

    # package otlphttp needs golang >= v1.20
    if [ "otlphttp" = $(basename $dirpath) ]; then
        if ! go version|grep -qE "go1.[2-9][0-9]"; then
          echo "ignore otlphttp as go version: $(go version)"
          continue 1
        fi
    fi

    # package otelmetric needs golang >= v1.20
    if [ "otelmetric" = $(basename $dirpath) ]; then
        if ! go version|grep -qE "go1.[2-9][0-9]"; then
          echo "ignore otelmetric as go version: $(go version)"
          continue 1
        fi
    fi

    cd $dirpath
    go mod tidy
    go build ./...
    # check coverage
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
