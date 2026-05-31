# gpt-private GoPay checkout sidecar

This package keeps only GPT-owned ChatGPT checkout / Stripe / Midtrans `snap_token` preparation used by `gpt-service`.

GoPay App account/device/user-state and Midtrans + GoPay linking/payment runtime are owned by the standalone `gopay-app` service.

## Entry point

- `cmd/gopay-payment-server`
- `proto/payment.proto`
