#!/bin/bash
set -ueo pipefail

eval `go env | sed -e 's/^/export /'`

if [ "${1:-}" = "watch" ]; then
  $GOPATH/bin/goconvey .
else
  go test -coverprofile coverage.out -json ./... | tee test-report.json
  go tool cover -html=coverage.out -o coverage.html
fi

