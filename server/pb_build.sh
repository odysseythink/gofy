#!/bin/sh

# request protobuf build
protoc -I=proto --go_out=proto --go_opt=paths=source_relative request/common.proto

# entities protobuf build
protoc -I=proto --go_out=proto --go_opt=paths=source_relative exceptions/exceptions.proto

# entities protobuf build
protoc -I=proto --go_out=proto --go_opt=paths=source_relative entities/feature.proto
protoc -I=proto --go_out=proto --go_opt=paths=source_relative entities/app.proto

# grpc protobuf build
protoc -I=proto --go_out=proto --go_opt=paths=source_relative   --go-grpc_out=proto --go-grpc_opt=paths=source_relative pbapi/link.proto
protoc -I=proto --go_out=proto --go_opt=paths=source_relative   --go-grpc_out=proto --go-grpc_opt=paths=source_relative pbapi/admin.proto
protoc -I=proto --go_out=proto --go_opt=paths=source_relative   --go-grpc_out=proto --go-grpc_opt=paths=source_relative pbapi/app.proto
protoc -I=proto --go_out=proto --go_opt=paths=source_relative   --go-grpc_out=proto --go-grpc_opt=paths=source_relative pbapi/tools.proto
protoc -I=proto --go_out=proto --go_opt=paths=source_relative   --go-grpc_out=proto --go-grpc_opt=paths=source_relative pbapi/plugins.proto