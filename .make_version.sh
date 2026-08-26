#!/usr/bin/env bash

# Function to run sed in-place with OS-specific options
sed_replace() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS - requires empty string after -i
        sed -i '' "$@"
    else
        # Linux/Windows Git Bash
        sed -i "$@"
    fi
}

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

# Parse the replace block from cmd/gf/go.work to keep the replace directives
# in sync with the go.work file. Each line inside the "replace (" block is
# expected to be in the form "module/path => ../../relative/path". The output
# is a list of "module=target" pairs, one per line.
parse_replace_modules() {
    local goWorkFile="cmd/gf/go.work"
    if [ ! -f "$goWorkFile" ]; then
        echo "Error: $goWorkFile not found, cannot sync replace directives"
        exit 1
    fi
    # Extract module and target paths between "replace (" and the closing ")".
    sed -n '/^replace (/,/^)/p' "$goWorkFile" \
        | grep -E '=>' \
        | sed -E 's/^[[:space:]]*([^[:space:]]+)[[:space:]]*=>[[:space:]]*([^[:space:]]+).*/\1=\2/'
}

# Collect the replace module=target pairs once at startup.
replacePairs=$(parse_replace_modules)

# check find command support or not
output=$(find "${workdir}" -name go.mod 2>&1)
if [[ $? -ne 0 ]]; then
    echo "Error: please use bash or zsh to run!"
    exit 1
fi

if [[ true ]]; then
    # Use sed to replace the version number in version.go
    sed_replace 's/VERSION = ".*"/VERSION = "'${newVersion}'"/' version.go

    # Use sed to replace the version number in README.MD
    sed_replace 's/version=[^"]*/version='${newVersion}'/' README.MD
    sed_replace 's/version=[^"]*/version='${newVersion}'/' README.zh_CN.MD
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

    # Add replace directive for local development.
    if [ $goModPath = "./cmd/gf" ]; then
        mv go.work go.work.version.bak
        # Add replace directives dynamically from cmd/gf/go.work so that the
        # script stays in sync with the replace block in that file.
        for pair in $replacePairs; do
            go mod edit -replace "$pair"
        done
    fi
    # Remove indirect dependencies
    sed_replace '/\/\/ indirect/d' go.mod
    go mod tidy
    # Remove toolchain line if exists
    sed_replace '/^toolchain/d' go.mod

    # Upgrading only GoFrame related libraries, sometimes even if a version number is specified,
    # it may not be possible to successfully upgrade. Please confirm before submitting the code
    go list -f "{{if and (not .Indirect) (not .Main)}}{{.Path}}@${newVersion}{{end}}" -m all | grep "^github.com/gogf/gf"
    go list -f "{{if and (not .Indirect) (not .Main)}}{{.Path}}@${newVersion}{{end}}" -m all | grep "^github.com/gogf/gf" | xargs -L1 go get -v
    # Remove indirect dependencies
    sed_replace '/\/\/ indirect/d' go.mod
    go mod tidy
    # Remove toolchain line if exists
    sed_replace '/^toolchain/d' go.mod
    if [ $goModPath = "./cmd/gf" ]; then
        # Drop the replace directives that were added above, keeping the
        # script in sync with the replace block in cmd/gf/go.work.
        for pair in $replacePairs; do
            go mod edit -dropreplace "${pair%%=*}"
        done
        mv go.work.version.bak go.work
    fi
    cd -
done

if [ -f "go.work.version.bak" ]; then
    mv go.work.version.bak go.work
    echo "Restore the go.work file"
fi