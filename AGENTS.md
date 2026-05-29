# AGENTS.md

本仓是 `gpt-private` 私有 GPT 能力仓。

- 本仓承载不适合放入公开 `gpt` 仓的 GPT 私有 provider、GoPay 运行时、私有动作实现、私有 workflow 和私有部署辅助脚本。
- 破坏性迁移按目标职责直接演进，不维护兼容层、旧入口或双轨逻辑。
- 可被公开服务长期复用的契约优先留在 `gpt` 或 `common-lib` 的 proto 真源；本仓只放私有实现、私有 provider raw shape、私有 workflow 和私有配置 schema。
- 不跨仓直接 import sibling 源码实现细节；需要组合时通过公开 plugin 注册接口、proto/gRPC、HTTP 边界、事件或部署期装载完成。
- GoPay sidecar 只作为 `gpt-service` 的私有 runtime 扩展运行；主编排能力仍由 `gpt` 仓的 orchestrator 通过 gRPC/HTTP 边界调用。
- 私有 action/config 通过 `gpt/pkg/gptplugin` 注册；不得 import `gpt/orchestrator/internal/...`。
- 运行时 secret、token、cookie、验证码、支付凭据、代理凭据和会话材料按敏感数据处理，不提交真实值。
- 不新增测试目录、测试夹具、覆盖率报告或 CI/CD 配置。
