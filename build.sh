#!/usr/bin/env bash

set -e

## required for glide
export GO15VENDOREXPERIMENT=1


glide install

## build it
go clean -r
go clean -i
cd kvipamdriver
go build --ldflags '-extldflags "-static"' ./...
