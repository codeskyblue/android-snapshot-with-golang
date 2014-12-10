#! /bin/sh
#
# build.sh
# Copyright (C) 2014 hzsunshx <hzsunshx@onlinegame-13-180>
#
# Distributed under terms of the MIT license.
#


CGO_ENABLED=1 GOOS=android GOARCH=arm go build
