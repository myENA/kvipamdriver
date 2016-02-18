#!/usr/bin/env bash

set -ex

[ -z "${1}" ] && DISTFILE="kvipamdriver.tgz" || DISTFILE="${1}"


BINNAME=kvipamdriver
BINPATH="./${BINNAME}"
UNLINK=false


mkdir -p tmp/$BINNAME/

cd tmp/$BINNAME/

cp ../../LICENSE .
cp ../../README.md .
cp ../../kvipamdriver/$BINNAME .
cd ..

tar -czhf ../"${DISTFILE}" $BINNAME

cd ..
rm -rf tmp
