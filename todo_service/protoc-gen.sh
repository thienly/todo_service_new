#!/bin/bash
protoc --proto_path=api/proto v1/todo.proto v1/user.proto  --go_out=:pb \
  --go-grpc_out=:pb  --grpc-gateway_out=:pb  --openapiv2_out=:swagger --openapiv2_opt use_go_templates=true
