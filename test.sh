#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2023, deadc0de6

set -e

# deps
echo "install deps..."
go install golang.org/x/lint/golint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

# linting
echo "go fmt..."
go fmt ./...
echo "golint..."
golint -set_exit_status ./...
echo "staticcheck..."
staticcheck ./...
echo "go vet..."
go vet ./...

# test scripts
echo "lint shell scripts"
find . -iname '*.sh' | while read -r script; do
  shellcheck "${script}"
done

# test python
echo "lint python scripts"
find . -iname '*.py' | while read -r script; do
  pylint -sn "${script}"
done

# compilation
echo "compiling..."
make clean
make

# diffs
echo "run tests..."
./tests/launcher.sh

echo "everything is OK!"
