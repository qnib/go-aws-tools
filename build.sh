#!/bin/bash

pushd cmd/query-regions
CGO_ENABLED=0 GOOS=linux go build -o query-regions_linux
go build -o query-regions_darwin
popd
