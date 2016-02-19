#!/usr/bin/env bash

set -e

## required for glide
export GO15VENDOREXPERIMENT=1


UPDATE=false
INSTALL=false

usage() {
    scriptname=`basename "$0"`
    echo
    echo Usage
    echo "$scriptname [-i] [-u]"
    echo
    echo "    -i is for install.  this runs glide install before build."
    echo
    echo "    -u is for update.  This runs a glide update before build."
}

while getopts :ui opt; do
    case "${opt}" in
        u)
            UPDATE=true
            ;;
        i)
            INSTALL=true
            ;;
        *)
            usage
            exit 1
            ;;
    esac
done

shift $((OPTIND-1))

if $UPDATE; then
    ## update glide
    printf "Refreshing vendor directory ... "
    glide -q update
elif $INSTALL; then
    ##install glide
    printf "Installing glide.lock"
    glide -q install
fi

## build it
go clean -r
go clean -i
cd kvipamdriver
CGO_ENABLED=0 go build  ./...
