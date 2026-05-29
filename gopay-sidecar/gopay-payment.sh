#!/usr/bin/env sh
set -eu
listen="${GOPAY_PAYMENT_LISTEN_ADDR:-:50054}"
exec /app/bin/gopay-payment --listen "$listen"
