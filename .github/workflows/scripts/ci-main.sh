#!/usr/bin/env bash

group=$1
coverage=$2

race=-race
if [ "$group" = "386" ]; then
    export GOARCH=386
    race=
    coverage=
fi

module_group() {
    case $1 in
        ./contrib/drivers/dm|./contrib/drivers/mssql|./contrib/drivers/oceanbase|./contrib/drivers/oracle|./contrib/drivers/sqlitecgo)
            echo "drivers-1" ;;
        ./contrib/drivers/*)
            echo "drivers-2" ;;
        *)
            echo "core" ;;
    esac
}

# go_version_at_least checks and returns whether the running go command satisfies the given
# version, like "1.26" or "1.26.0".
go_version_at_least() {
    local required=$1
    local current
    current=$(go env GOVERSION | sed 's/^go//; s/[^0-9.].*$//')
    # Pad the missing minor/patch segments so that "1.26" and "1.26.0" compare equal.
    while [ "$(echo "$required" | tr -cd '.' | wc -c)" -lt 2 ]; do
        required="$required.0"
    done
    while [ "$(echo "$current" | tr -cd '.' | wc -c)" -lt 2 ]; do
        current="$current.0"
    done
    [ "$(printf '%s\n%s\n' "$required" "$current" | sort -V | head -1)" = "$required" ]
}

# module_go_version retrieves the `go` directive of the go.mod in given directory.
module_go_version() {
    sed -n 's/^go[[:space:]]\{1,\}\([0-9][0-9.]*\).*/\1/p' "$1/go.mod" | head -1
}

# find all path that contains go.mod.
for file in `find . -name go.mod`; do
    dirpath=$(dirname $file)
    echo $dirpath
    
    # package kubecm was moved to sub ci procedure.
    if [ "kubecm" = $(basename $dirpath) ]; then
        continue 1
    fi
    
    # examples directory was moved to sub ci procedure.
    if [[ $dirpath =~ "/examples/" ]]; then
        continue 1
    fi
    
    if [[ $file =~ "/testdata/" ]]; then
        echo "ignore testdata path $file"
        continue 1
    fi
    
    if [ "$group" = "386" ]; then
        if [ "$dirpath" != "." ]; then
            continue 1
        fi
    elif [ "$(module_group $dirpath)" != "$group" ]; then
        continue 1
    fi
    
    # Check if it's a contrib directory
    if [[ $dirpath =~ "/contrib/" ]]; then
        # Check if go version meets the requirement
        if ! go version | grep -qE "go${LATEST_GO_VERSION}"; then
            echo "ignore path $dirpath as go version is not ${LATEST_GO_VERSION}: $(go version)"
            # clean docker containers and images to free disk space
            # bash .github/workflows/scripts/ci-main-clean.sh "$dirpath"
            continue 1
        fi
    fi

    # A module may require a newer go version than an older matrix entry provides, for
    # example `cmd/gf` requires go 1.26.0 because of its `golang.org/x/tools` dependency.
    # Skip it on the older entries instead of failing the whole matrix leg.
    moduleGoVersion=$(module_go_version "$dirpath")
    if [ -n "$moduleGoVersion" ] && ! go_version_at_least "$moduleGoVersion"; then
        echo "ignore path $dirpath as it requires go >= ${moduleGoVersion}: $(go version)"
        continue 1
    fi
    
    # if [[ $dirpath = "." ]]; then
    #     # No space left on device error sometimes occurs in CI pipelines, so clean the cache before tests.
    #     go clean -cache
    # fi
    
    cd $dirpath
    go mod tidy
    go build ./...
    
    # test with coverage
    if [ "${coverage}" = "coverage" ]; then
        go test ./... -count=1 $race -coverprofile=coverage.out -covermode=atomic -coverpkg=./...,github.com/gogf/gf/... || exit 1
        
        if grep -q "/gogf/gf/.*/v2" go.mod; then
            sed -i "s/gogf\/gf\(\/.*\)\/v2/gogf\/gf\/v2\1/g" coverage.out
        fi
    else
        go test ./... -count=1 $race || exit 1
    fi
    
    cd -
    # Docker cleanup is intentionally disabled. The services are started once per job
    # in the "Start Services" step and shared by every module of the group, so tearing
    # them down after a single module would break the modules that follow. The master
    # runs complete the drivers groups without it, so there is no disk pressure to
    # relieve here either.
    # bash .github/workflows/scripts/ci-main-clean.sh "$dirpath"
done
