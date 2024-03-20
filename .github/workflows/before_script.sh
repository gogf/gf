#!/usr/bin/env bash

find . -name "*.go" | xargs gofmt -w
git diff --name-only --exit-code || if [ $? != 0 ]; then
    echo "Notice: gofmt checks have failed, please gofmt before pr." && exit 1;
fi
echo "gofmt checks have passed."

find . -name "*_test.go" -print0 | while IFS= read -r -d '' file; do
    awk '/func Test[[:upper:]]/ { print $2 }' "$file" | while read -r funcName; do
        if [[ ! $funcName =~ ^Test[[:upper:]] ]]; then
            echo "Notice: Func name $funcName in file $file checks have failed, please check that it is upper camel case before pr."
        fi
    done
done

echo "Func name of unit test checks have passed."


sudo echo "127.0.0.1   local" | sudo tee -a /etc/hosts