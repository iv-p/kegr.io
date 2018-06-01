#!/usr/bin/env sh
set -ex

dep ensure  -vendor-only -v

exec $@