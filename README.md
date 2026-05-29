# GPT Private

本仓承载 GPT 私有 provider、私有动作元数据、私有 workflow 和 GoPay sidecar runtime。

## 目录

- `plugins/`：通过 `gpt/pkg/gptplugin` 注册私有 action/config/workflow 元数据。
- `gopay/`：GoPay App / payment gRPC sidecar 源码、私有 proto 真源、协议 client、Android OTP relay 和 GoPay 参考资料。
- `gopay-sidecar/`：部署到 `gpt-service` 私有镜像中的 sidecar 启动脚本。
- `workflows/`：私有 n8n workflow 定义。

## 生成

GoPay sidecar proto 只在本仓生成：

```bash
(cd gopay && ./scripts/generate-proto.sh)
```

## 约束

- 不 import `gpt/orchestrator/internal/...`。
- GoPay 具体 action id、step id、provider 配置和 workflow 均由本仓持有。
- `gpt` 缺少本仓时必须能正常编译基础能力。
