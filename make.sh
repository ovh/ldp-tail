#!/bin/bash

BUILD_VERSION=${BUILD_VERSION:-"dev"}
BUILD_TIME=`date +"%Y-%m-%dT%H:%M:%S%z"`
GIT_TAG=`git describe --always --dirty=-wip`

go get -v -d
go build -ldflags "-X stash.ovh.net/framework/build.Version=$BUILD_VERSION -X stash.ovh.net/framework/build.Time=$BUILD_TIME -X stash.ovh.net/framework/build.GitTag=$GIT_TAG"
