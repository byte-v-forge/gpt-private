//go:build private_plugins

package app

import (
	"context"

	"github.com/byte-v-forge/common-lib/envx"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"orchestrator/internal/activities"
	"orchestrator/internal/api"
	"orchestrator/internal/contracts"
	"orchestrator/pb"
)

func configurePrivateDependencies(ctx context.Context, cfg orchestratorConfig, deps *orchestratorDependencies, _ redis.Cmdable) error {
	if !deps.goPayActionsEnabled() {
		activities.ConfigureGoPayDependencies(nil)
		return nil
	}
	conn, err := newGRPCClientConn(
		"GoPay app channel",
		envx.StringDefault("GPT_PRIVATE_APP_ADDR", envx.StringDefault("GPT_PRIVATE_APP_INTERNAL_ADDR", "gopay-app:50051")),
		grpc.WithDefaultServiceConfig(gopayAppGRPCRetryServiceConfig()),
	)
	if err != nil {
		return err
	}
	deps.addCloser(conn.Close)
	client := pb.NewGopayAppServiceClient(conn)
	api.ConfigureGoPayClient(client)
	activities.ConfigureGoPayDependencies(client)
	return nil
}

func (d *orchestratorDependencies) goPayActionsEnabled() bool {
	return d.hasAnyAction(
		contracts.ActionGoPayApp,
		contracts.ActionGoPayPayment,
		contracts.ActionGoPayQRISPaymentActivate,
		contracts.ActionGoPayWAPayment,
		contracts.ActionGoPayPaymentRebind,
	)
}

func gopayAppGRPCRetryServiceConfig() string {
	return `{
		"methodConfig": [{
			"name": [{"service": "gopay_app.GopayAppService"}],
			"retryPolicy": {
				"MaxAttempts": 3,
				"InitialBackoff": "0.3s",
				"MaxBackoff": "2s",
				"BackoffMultiplier": 2,
				"RetryableStatusCodes": ["UNAVAILABLE"]
			}
		}]
	}`
}
