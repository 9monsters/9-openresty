#!/bin/sh
set -eu
export PATH="/usr/local/bin:$HOME/go/bin:$PATH"

echo "protoc: $(which protoc) $(protoc --version)"
echo "protoc-gen-go: $(which protoc-gen-go 2>/dev/null || echo 'installing...')"
echo "protoc-gen-go-grpc: $(which protoc-gen-go-grpc 2>/dev/null || echo 'installing...')"

if ! command -v protoc-gen-go >/dev/null 2>&1; then
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi
if ! command -v protoc-gen-go-grpc >/dev/null 2>&1; then
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

cd /Users/wangjia/projects/9-openresty

protoc \
  --proto_path=proto \
  --go_out=. --go_opt=module=openresty-admin \
  --go-grpc_out=. --go-grpc_opt=module=openresty-admin \
  proto/openresty/admin/v1/services.proto

echo "=== Generated files ==="
ls -la gen/go/openresty/admin/v1/