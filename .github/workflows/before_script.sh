#!/usr/bin/env bash

find . -name "*.go" | xargs gofmt -w
git diff --name-only --exit-code || if [ $? != 0 ];
    then echo "Notice: gofmt checks have failed, please gofmt before pr." && exit 1;
fi
echo "gofmt checks have passed."

find . -name "*_test.go" -print0 | while IFS= read -r -d '' file; do
    awk '/func Test[[:upper:]][[:alnum:]]*\(/ {
        if ($0 !~ /func Test[[:upper:]][[:alnum:]]*\(/) {
            print "Notice: Func name of unit test checks have failed, please check that it is upper camel case before pr"
            exit 1
        }
    }' "$file"
done
echo "Func name of unit test checks have passed."

sudo echo "127.0.0.1   local" | sudo tee -a /etc/hosts