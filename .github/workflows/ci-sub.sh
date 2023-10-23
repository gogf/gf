#!/usr/bin/env bash

coverage=$1

# find all path that contains go.mod.
for file in `find . -name go.mod`; do
    dirpath=$(dirname $file)
    echo $dirpath

    # package kuhecm needs golang >= v1.19
    if [ "kubecm" = $(basename $dirpath) ]; then
        if ! go version|grep -qE "go1.19|go1.[2-9][0-9]"; then
          echo "ignore kubecm as go version: $(go version)"
          continue 1
        fi
    else
      continue 1
    fi

    cd $dirpath

    go mod tidy
    go build ./...
    go test ./... -race || exit 1

    cd -
done
