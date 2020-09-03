#!/bin/bash
set -ueo pipefail

eval `go env | sed -e 's/^/export /'`

if [ "${1:-}" = "watch" ]; then
  $GOPATH/bin/goconvey .
else
  go test ./...
fi

