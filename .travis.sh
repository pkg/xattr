#!/bin/sh -e

echo "Building for Linux..."
GOOS=linux   GOARCH=amd64 go build
echo "Building for Darwin..."
GOOS=darwin  GOARCH=amd64 go build
echo "Building for FreeBSD..."
GOOS=freebsd GOARCH=amd64 go build

echo "Running tests.."
uname -a
go vet
go test -v -race
