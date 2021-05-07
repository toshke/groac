#!/bin/bash

set -eo pipefail

BUILD_LOC="$PROJECT_ROOT/build"
rm -rf "$BUILD_LOC" && mkdir -p "$BUILD_LOC"
go build -o "${BUILD_LOC}/"
