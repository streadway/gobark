#!/bin/bash
set +e

go build .
mkdir -p usr/bin
cp gobark usr/bin

ver=$(git describe | sed s/^v// | tr - + | tr -d '\n')

git diff --quiet HEAD || ver="$ver+mod"

fpm -m "Sean Treadway <sean@soundcloud.com>" \
  -s dir \
  -t deb \
  -v $ver \
  -n gobark \
  usr/bin/

rm -rf usr
