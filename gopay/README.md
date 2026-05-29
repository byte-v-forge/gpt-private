# GoPay

本目录承载 GoPay 协议实现和通道运行时代码。

## 目录

- `protocol/`：GoPay App 和钱包侧支付调用使用的 Go 协议包。
- `proto/`：GoPay sidecar 私有 proto 真源。
- `whatsapp-relay/`：GoPay OTP Android 通知转发器。
- `channels/payment/`：GoPay 支付链路参考资料与外部许可文件。

## 检查

构建检查只在远程构建环境执行；本机提交前仅运行 `gofmt` 和静态差异检查。
