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

workdir=.
echo "Prepare to tidy all go.mod files in the ${workdir} directory"

# check find command support or not
output=$(find "${workdir}" -name go.mod 2>&1)
if [[ $? -ne 0 ]]; then
    echo "Error: please use bash or zsh to run!"
    exit 1
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
    # Remove indirect dependencies
    sed_replace '/\/\/ indirect/d' go.mod
    go mod tidy
    # Remove toolchain line if exists
    sed_replace '/^toolchain/d' go.mod
    cd - > /dev/null
done
