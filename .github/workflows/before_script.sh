#!/usr/bin/env bash

find . -name "*.go" | xargs gofmt -w
git diff --name-only --exit-code || if [ $? != 0 ]; then echo "Notice: gofmt check failed,please gofmt before pr." && exit 1; fi
echo "gofmt check pass."
sudo echo "127.0.0.1   local" | sudo tee -a /etc/hosts