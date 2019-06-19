#!/bin/bash

cd "${GOPATH}/src/github.com/pibigstar/go-todo"

find . -name "*.go" -not -path "./third/*" | xargs gofmt -w

git diff --exit-code
