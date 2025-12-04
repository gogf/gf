#!/usr/bin/env bash

# Get the absolute path to the repository root
repo_root=$(pwd)
workdir=$repo_root/examples

echo "Prepare to process go.mod files in the ${workdir} directory"

# Check if examples directory exists
if [ ! -d "${workdir}" ]; then
    echo "Error: examples directory not found at ${workdir}"
    exit 1
fi

# Check if find command is available
if ! command -v find &> /dev/null; then
    echo "Error: find command not found!"
    exit 1
fi

for file in `find ${workdir} -name go.mod`; do
    goModPath=$(dirname $file)
    echo ""
    echo "Processing dir: $goModPath"

    # Calculate relative path to root
    # First get the relative path from go.mod to repo root
    relativePath=""
    current="$goModPath"
    while [ "$current" != "$repo_root" ]; do
        relativePath="../$relativePath"
        current=$(dirname "$current")
    done
    relativePath=${relativePath%/}  # Remove trailing slash
    echo "Relative path to root: $relativePath"

    # Get all github.com/gogf/gf dependencies
    # Use awk to get package names without version numbers
    dependencies=$(awk '/^[[:space:]]*github\.com\/gogf\/gf\// {print $1}' "$file" | sort -u)
    
    if [ -n "$dependencies" ]; then
        echo "Found GoFrame dependencies:"
        echo "$dependencies"
        echo "Adding replace directives..."

        # Create temporary file
        temp_file="${file}.tmp"
        # Remove existing replace directives and copy to temp file
        sed '/^replace.*github\.com\/gogf\/gf.*/d' "$file" > "$temp_file"
        
        # Add new replace block
        echo "" >> "$temp_file"
        echo "replace (" >> "$temp_file"
        
        while IFS= read -r dep; do
            # Skip empty lines
            [ -z "$dep" ] && continue
            
            # Calculate the relative path for the replacement
            if [[ "$dep" == "github.com/gogf/gf/v2" ]]; then
                replacement="$relativePath"
            else
                # Extract the path after v2 and remove trailing version
                subpath=$(echo "$dep" | sed -E 's/github\.com\/gogf\/gf\/(contrib\/[^/]+\/[^/]+)\/v2.*/\1/')
                replacement="$relativePath/$subpath"
            fi
            
            echo "    $dep => $replacement/" >> "$temp_file"
        done <<< "$dependencies"
        
        echo ")" >> "$temp_file"
        
        # Replace original file with temporary file
        mv "$temp_file" "$file"
        echo "Replace directives added to $file"
    else
        echo "No GoFrame dependencies found in $file"
    fi
done

echo "\nAll go.mod files have been processed successfully."
