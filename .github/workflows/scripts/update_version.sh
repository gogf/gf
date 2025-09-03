#!/usr/bin/env bash

# Check if the number of parameters is 2
if [ $# -ne 2 ]; then
    echo "Invalid parameters, please execute in format: version.sh [directory] [version]"
    echo "Example: version.sh ./contrib v1.0.0"
    exit 1
fi

# Check if the first parameter is a directory and exists
if [ ! -d "$1" ]; then
    echo "Error: Directory does not exist"
    exit 1
fi

# Check if the second parameter starts with 'v'
if [[ "$2" != v* ]]; then
    echo "Error: Version number does not start with 'v'"
    exit 1
fi

workdir=$1
newVersion=$2
echo "Preparing to replace version numbers in all go.mod files under ${workdir} directory to ${newVersion}"


# Check if file exists
if [ -f "go.work" ]; then
    # File exists, rename it
    mv go.work go.work.${newVersion}
    echo "Backup go.work file to avoid affecting the upgrade"
fi

for file in `find ${workdir} -name go.mod`; do
    goModPath=$(dirname $file)
    echo ""
    echo "processing dir: $goModPath"
    cd $goModPath
    go mod tidy
    go list -f "{{if and (not .Indirect) (not .Main)}}{{.Path}}@${newVersion}{{end}}" -m all | grep "^github.com/gogf/gf"
    go list -f "{{if and (not .Indirect) (not .Main)}}{{.Path}}@${newVersion}{{end}}" -m all | grep "^github.com/gogf/gf" | xargs -L1 go get -v 
    go mod tidy
    cd -
done

if [ -f "go.work.${newVersion}" ]; then
    # File exists, rename it back
    mv go.work.${newVersion} go.work
    echo "Restore go.work file"
fi
