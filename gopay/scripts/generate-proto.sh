#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "${ROOT}/pb"
rm -f "${ROOT}/pb"/*.pb.go "${ROOT}/pb"/*_grpc.pb.go
protoc -I "${ROOT}/proto" \
  --go_out="${ROOT}/pb" \
  --go-grpc_out="${ROOT}/pb" \
  "${ROOT}/proto/gopay_app.proto" \
  "${ROOT}/proto/payment.proto"
