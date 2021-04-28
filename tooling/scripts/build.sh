#!/bin/bash

set -eo pipefail

cd "$PROJECT_ROOT/cmd"
BUILD_LOC="$PROJECT_ROOT/build"
rm -rf "$BUILD_LOC" && mkdir -p "$BUILD_LOC"

for cmd in *; do
    cd "${PROJECT_ROOT}/cmd/${cmd}"
    go build -o "${BUILD_LOC}/${cmd}"
done