#!/bin/bash
cd ..
mv LightRPC/ rpc/
go get go.uber.org/zap
go get github.com/BurntSushi/toml
go get github.com/modood/table
go get github.com/golang/protobuf/proto

