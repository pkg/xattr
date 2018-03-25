#!/bin/sh -e

echo "Building for Linux..."
GOOS=linux   go build
GOOS=linux   go vet

echo "Building for Darwin..."
GOOS=darwin  go build
GOOS=darwin  go vet

echo "Building for FreeBSD..."
GOOS=freebsd go build
GOOS=freebsd go vet

echo "Running tests..."
go test -v -race
