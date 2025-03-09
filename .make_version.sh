#!/usr/bin/env bash
if [ $# -ne 2 ]; then
    echo "Parameter exception, please execute in the format of $0 [directory] [version number]"
    echo "PS：$0 ./ v2.4.0"
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

workdir=.
newVersion=$2
echo "Prepare to replace the GoFrame library version numbers in all go.mod files in the ${workdir} directory with ${newVersion}"

# check find command support or not
output=$(find "${workdir}" -name go.mod 2>&1)
if [[ $? -ne 0 ]]; then
    echo "Error: please use bash or zsh to run!"
    exit 1
fi

if [[ true ]]; then
    # Use sed to replace the version number in version.go
    sed -i '' 's/VERSION = ".*"/VERSION = "'${newVersion}'"/' version.go

    # Use sed to replace the version number in README.MD
    sed -i '' 's/version=[^"]*/version='${newVersion}'/' README.MD
fi

if [ -f "go.work" ]; then
    mv go.work go.work.version.bak
    echo "Back up the go.work file to avoid affecting the upgrade"
fi

for file in `find ${workdir} -name go.mod`; do
    goModPath=$(dirname $file)
    echo ""
    echo "processing dir: $goModPath"

    if [[ $goModPath =~ "/testdata/" ]]; then
        echo "ignore testdata path $goModPath"
        continue 1
    fi

    if [[ $goModPath =~ "/examples/" ]]; then
        echo "ignore examples path $goModPath"
        continue 1
    fi

    cd $goModPath
    if [ $goModPath = "./cmd/gf" ]; then
        mv go.work go.work.version.bak
        go mod edit -replace github.com/gogf/gf/v2=../../
        go mod edit -replace github.com/gogf/gf/contrib/drivers/clickhouse/v2=../../contrib/drivers/clickhouse
        go mod edit -replace github.com/gogf/gf/contrib/drivers/mssql/v2=../../contrib/drivers/mssql
        go mod edit -replace github.com/gogf/gf/contrib/drivers/mysql/v2=../../contrib/drivers/mysql
        go mod edit -replace github.com/gogf/gf/contrib/drivers/oracle/v2=../../contrib/drivers/oracle
        go mod edit -replace github.com/gogf/gf/contrib/drivers/pgsql/v2=../../contrib/drivers/pgsql
        go mod edit -replace github.com/gogf/gf/contrib/drivers/sqlite/v2=../../contrib/drivers/sqlite
    fi
    go mod tidy
    # Remove toolchain line if exists
    sed -i '' '/^toolchain/d' go.mod

    # Upgrading only GoFrame related libraries, sometimes even if a version number is specified, 
    # it may not be possible to successfully upgrade. Please confirm before submitting the code
    go list -f "{{if and (not .Indirect) (not .Main)}}{{.Path}}@${newVersion}{{end}}" -m all | grep "^github.com/gogf/gf"
    go list -f "{{if and (not .Indirect) (not .Main)}}{{.Path}}@${newVersion}{{end}}" -m all | grep "^github.com/gogf/gf" | xargs -L1 go get -v 
    go mod tidy
    # Remove toolchain line if exists
    sed -i '' '/^toolchain/d' go.mod
    if [ $goModPath = "./cmd/gf" ]; then
        go mod edit -dropreplace github.com/gogf/gf/v2
        go mod edit -dropreplace github.com/gogf/gf/contrib/drivers/clickhouse/v2
        go mod edit -dropreplace github.com/gogf/gf/contrib/drivers/mssql/v2
        go mod edit -dropreplace github.com/gogf/gf/contrib/drivers/mysql/v2
        go mod edit -dropreplace github.com/gogf/gf/contrib/drivers/oracle/v2
        go mod edit -dropreplace github.com/gogf/gf/contrib/drivers/pgsql/v2
        go mod edit -dropreplace github.com/gogf/gf/contrib/drivers/sqlite/v2
        mv go.work.version.bak go.work
    fi
    cd -
done

if [ -f "go.work.version.bak" ]; then
    mv go.work.version.bak go.work
    echo "Restore the go.work file"
fi