#!/usr/bin/env bash

coverage=$1

# update code of submodules
git clone https://github.com/gogf/examples

# update go.mod in examples directory to replace github.com/gogf/gf packages with local directory
bash .github/workflows/scripts/replace_examples_gomod.sh

# Function to compare version numbers
version_compare() {
    local ver1=$1
    local ver2=$2
    
    # Remove 'go' prefix and 'v' if present
    ver1=$(echo "$ver1" | sed 's/^go//; s/^v//')
    ver2=$(echo "$ver2" | sed 's/^go//; s/^v//')
    
    # Split versions into major.minor format
    local major1=$(echo "$ver1" | cut -d. -f1)
    local minor1=$(echo "$ver1" | cut -d. -f2)
    local major2=$(echo "$ver2" | cut -d. -f1)
    local minor2=$(echo "$ver2" | cut -d. -f2)
    
    # Compare versions: return 0 if ver1 <= ver2, 1 otherwise
    if [ "$major1" -lt "$major2" ]; then
        return 0
    elif [ "$major1" -eq "$major2" ] && [ "$minor1" -le "$minor2" ]; then
        return 0
    else
        return 1
    fi
}

# Get current Go version
current_go_version=$(go version | grep -oE 'go[0-9]+\.[0-9]+')

# find all path that contains go.mod.
for file in `find . -name go.mod`; do
    dirpath=$(dirname $file)
    echo "Processing: $dirpath"

    # Only process examples and kubecm directories  

    # Process examples directory (only build, no tests)
    if [[ $dirpath =~ "/examples/" ]]; then
        echo "  the examples directory only needs to be built, not unit tests."
        cd $dirpath
        go mod tidy
        go build ./...
        cd -
        continue 1
    fi
    
    # Process kubecm directory
    if [ "kubecm" != $(basename $dirpath) ]; then
        echo "  Skipping: not kubecm directory"
        continue
    fi

    cd $dirpath

    # Read Go version requirement from go.mod
    if [ -f "go.mod" ]; then
        go_mod_version=$(grep '^go ' go.mod | awk '{print $2}' | head -1)
        
        if [ -n "$go_mod_version" ]; then
            echo "  go.mod requires: go$go_mod_version"
            echo "  current version: $current_go_version"
            
            # Check if go.mod version requirement is satisfied by current Go version
            if version_compare "$go_mod_version" "$current_go_version"; then
                echo "  ✓ Version requirement satisfied, proceeding with build and test"
                
                go mod tidy
                go build ./...
                go test ./... -race || exit 1
            else
                echo "  ✗ Current Go version ($current_go_version) does not meet requirement (go$go_mod_version), skipping"
            fi
        fi
    fi

    cd -
done
