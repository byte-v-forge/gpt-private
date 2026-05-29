#!/usr/bin/env sh
set -eu
export GOPAY_APP_PORT="${GOPAY_APP_PORT:-50060}"
exec /app/bin/gopay-app
