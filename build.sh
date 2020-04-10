#!/usr/bin/env bash
GO_LOCATION=$(which go)
GOROOT=$(echo ${GO_LOCATION%/bin/go})

protoc --go_out=plugins=grpc:. protobuf/*.proto
protoc --go_out=plugins=grpc:. protobuf/jsonRPC/*.proto
go build -o ./bin/dapi .