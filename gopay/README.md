# gpt-private GoPay checkout sidecar

This package implements the standalone GPT checkout `payment.PaymentService` used by `gpt-service` for ChatGPT checkout / Stripe / Midtrans `snap_token` preparation and plan probing.

GoPay App account/device/user-state and Midtrans + GoPay linking/payment runtime are owned by the standalone `gopay-app` service.

## Entry point

- `cmd/gopay-payment-server`
- `Dockerfile`

## Contract

`payment.PaymentService` is sourced from `gpt/proto/payment.proto`. This package only generates local Go bindings into `pb/` during build or via:

```sh
./scripts/generate-proto.sh
```
