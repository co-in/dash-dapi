#!/usr/bin/env bash
GO_LOCATION=$(which go)
GOROOT=$(echo ${GO_LOCATION%/bin/go})

protoc --go_out=plugins=grpc:. evo/protobuf/*.proto
go build -o ./bin/dapi .