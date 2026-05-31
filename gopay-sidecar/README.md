# GoPay checkout sidecar

Only GPT checkout / Stripe / Midtrans `snap_token` preparation remains as a `gpt-service` sidecar.
GoPay App account/device/user-state and GoPay payment runtime are owned by the standalone `gopay-app` service.
