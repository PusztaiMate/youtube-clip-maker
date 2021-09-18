#!/usr/bin/env bash

protoc --proto_path=. --go_out=./backend/go --go_opt=paths=source_relative --go-grpc_out=./backend/go  --go-grpc_opt=paths=source_relative clips/clips.proto
