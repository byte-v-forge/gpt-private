#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SOURCE_ROOT="${SOURCE_ROOT:-$(cd "${ROOT}/../.." && pwd)}"
COMMON_PROTO_DIR="${COMMON_PROTO_DIR:-${SOURCE_ROOT}/common-lib/proto}"
GPT_PROTO_DIR="${GPT_PROTO_DIR:-${SOURCE_ROOT}/gpt/proto}"
mkdir -p "${ROOT}/pb"
rm -f "${ROOT}/pb"/*.pb.go "${ROOT}/pb"/*_grpc.pb.go
if [[ ! -f "${GPT_PROTO_DIR}/payment.proto" ]]; then
  printf 'required GPT payment proto not found under: %s\n' "${GPT_PROTO_DIR}" >&2
  exit 1
fi
protoc -I "${GPT_PROTO_DIR}" -I "${COMMON_PROTO_DIR}" \
  --go_out="${ROOT}/pb" \
  --go-grpc_out="${ROOT}/pb" \
  "${GPT_PROTO_DIR}/payment.proto"
