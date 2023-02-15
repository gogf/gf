#!/bin/bash
# update go.mod
echo "start search go.mod file"
for file in $(find . -name go.mod); do
    tempPath=$(dirname $file)
    echo "Processing dir: ${tempPath}"
    if [[ ${tempPath} =~ ".git" || ${tempPath} == "." ]] ; then
        echo "Skip path"
    elif [[ ${tempPath} =~ "./cmd/gf" || ${tempPath} =~ "./example" ]] ; then
        cd ${tempPath}
        go get -u -v ./...
        go mod tidy
        cd -
    else
        # just update gf
        cd ${tempPath}
        go get -u -v github.com/gogf/gf/v2
        go mod tidy
        cd -
    fi
done