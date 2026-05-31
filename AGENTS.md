# AGENTS.md

本仓是 `gpt-private` 私有 GPT 能力仓。

- 本仓承载不适合放入公开 `gpt` 仓的 GPT 私有 provider、ChatGPT checkout/Stripe 私有准备逻辑、私有动作实现、私有 workflow 和私有部署辅助脚本。
- 破坏性迁移按目标职责直接演进，不维护兼容层、旧入口或双轨逻辑。
- 可被公开服务长期复用的契约优先留在 `gpt` 或 `common-lib` 的 proto 真源；本仓只放私有实现、私有 provider raw shape、私有 workflow 和私有配置 schema。
- 不跨仓直接 import sibling 源码实现细节；需要组合时通过公开 plugin 注册接口、proto/gRPC、HTTP 边界、事件或部署期装载完成。
- GoPay checkout sidecar 只负责 GPT checkout/Stripe/snap_token 准备；GoPay App 与 payment runtime 归属 `gopay-app`。
- 私有 action/config 通过 `gpt/pkg/gptplugin` 注册；不得 import `gpt/orchestrator/internal/...`。
- 运行时 secret、token、cookie、验证码、支付凭据、代理凭据和会话材料按敏感数据处理，不提交真实值。
- Linter 检查必须达到 0 error / 0 warning；禁止通过修改或放宽 linter 配置、降低规则级别、删除规则、添加 ignore/disable/nolint/ts-ignore/eslint-disable/biome-ignore/prettier-ignore 等方式绕过问题，只能按 linter 规则修复源码、类型、格式或依赖边界。
- 不新增 CI/CD 配置。
