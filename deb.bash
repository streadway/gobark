#!/bin/bash
set +e

go build .
mkdir -p usr/bin
cp gobark usr/bin

ver=$(git describe | sed s/^v// | tr - + | tr -d '\n')
author=$(git log -1 $ver | grep 'Author: ' | sed 's/^Author: //')

fpm -m "$author" -v $ver -s dir -t deb -n gobark usr/bin/

rm -rf usr
