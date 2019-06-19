#!/bin/bash

cd "${GOPATH}/src/github.com/gogf/gf"

find . -name "*.go" -not -path "./third/*" | xargs gofmt -w

git diff --exit-code
