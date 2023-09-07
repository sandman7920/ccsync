#!/bin/bash

source "$HOME/cross-linaro/linaro7.5_env"

SELF="$(readlink -f "$0")"
HERE="${SELF%/*}"

cd "$HERE"

#VERBOSE="-v"
#GCC_VERBOSE="-x"
EXTRA="-ldflags=-s"
LTO="-flto"

export CGO_CFLAGS="-O2 -g -mtune=cortex-a9 $LTO"
export CGO_LDFLAGS="-O2 -g $LTO"
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=arm
export GOARM=7 
go build $VERBOSE $EXTRA $GCC_VERBOSE -o ./build && mv ./build/ccsync kual/extensions/ccsync/
