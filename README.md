# gpt-private

本仓承载 GPT 私有 provider、私有动作元数据、私有 workflow，以及 GPT checkout/Stripe/Midtrans `snap_token` 准备 sidecar 源码。

## 目录

- `plugins/`：通过 `gpt/pkg/gptplugin` 注册私有 action/config/workflow 元数据。
- `proto/`：本仓私有 action 元数据的 proto 真源。
- `gopay/`：GPT checkout sidecar 源码、协议 client 与准备逻辑；服务契约来自 `gpt/proto/payment.proto`。
- `gopay-sidecar/`：GoPay checkout sidecar 本地/镜像内启动脚本。
- `workflows/`：私有 n8n workflow JSON/catalog。
- `webui/`：私有 dashboard 扩展源码与 proto 生成脚本。

## 边界

- 本仓不再提供 `gpt-service` 镜像 overlay、不嵌入 `gpt/orchestrator`，也不承载 GoPay App/payment runtime 迁移。
- GoPay checkout sidecar 只负责 GPT checkout/Stripe/Midtrans `snap_token` 准备；GoPay App 与 payment runtime 归属 `gopay-app`。
- 私有 action/config 通过 `gpt/pkg/gptplugin` 注册，不 import `gpt/orchestrator/internal/...`。
- 私有 workflow 通过公开 GPT/gopay-app HTTP 或 gRPC 边界集成。

## Proto 生成

GoPay checkout sidecar proto 只在本仓生成：

```sh
cd gopay
./scripts/generate-proto.sh
```

WebUI 扩展 proto 类型通过 `webui/scripts/generate-proto.sh` 从 `gpt/proto`、`gopay-app/proto`、`common-lib/proto` 与本仓 `proto/` 生成。
