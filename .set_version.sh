#!/usr/bin/env bash

if [ $# -ne 2 ]; then
    echo "Parameter exception, please execute in the format of $0 [directory] [version number]"
    echo "PS：$0 ./contrib v1.0.0"
    exit 1
fi

if [ ! -d "$1" ]; then
    echo "Error: Directory does not exist"
    exit 1
fi

if [[ "$2" != v* ]]; then
    echo "Error: Version number must start with v"
    exit 1
fi

workdir=$1
newVersion=$2
echo "Prepare to replace the GF library version numbers in all go.mod files in the ${workdir} directory with ${newVersion}"

if [[ ${workdir} == ./contrib ]]; then
    echo "package gf" > version.go
    echo "" >> version.go
    echo "const (" >> version.go
    echo -e "\t// VERSION is the current GoFrame version." >> version.go
    echo -e "\tVERSION = \"${newVersion}\"" >> version.go
    echo ")" >> version.go
fi

if [ -f "go.work" ]; then
    mv go.work go.work.version.bak
    echo "Back up the go.work file to avoid affecting the upgrade"
fi

for file in `find ${workdir} -name go.mod`; do
    goModPath=$(dirname $file)
    echo ""
    echo "processing dir: $goModPath"
    cd $goModPath
    go mod tidy
    # Upgrading only GF related libraries, sometimes even if a version number is specified, it may not be possible to successfully upgrade. Please confirm before submitting the code
    go list -f "{{if and (not .Indirect) (not .Main)}}{{.Path}}@${newVersion}{{end}}" -m all | grep "^github.com/gogf/gf"
    go list -f "{{if and (not .Indirect) (not .Main)}}{{.Path}}@${newVersion}{{end}}" -m all | grep "^github.com/gogf/gf" | xargs -L1 go get -v 
    go mod tidy
    cd -
done

if [ -f "go.work.version.bak" ]; then
    mv go.work.version.bak go.work
    echo "Restore the go.work file"
fi