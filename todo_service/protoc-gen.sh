#!/bin/bash
rm -rf swagger/v1
mkdir swagger/v1
protoc --proto_path=pb todo.proto user.proto login.proto  --go_out=:pb \
  --go-grpc_out=:pb  --grpc-gateway_out=:pb  --openapiv2_out=:swagger/v1 --openapiv2_opt use_go_templates=true,allow_merge=true,merge_file_name=todo
