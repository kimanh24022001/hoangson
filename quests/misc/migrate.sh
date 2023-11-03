#!/usr/bin/bash

pushd cmd/reflect
go run main.go
popd


pushd cmd/migrate
go run main.go
popd
