//go:build private_plugins

package app

import (
	"google.golang.org/grpc"

	"orchestrator/internal/api"
	"orchestrator/pb"
)

func registerPrivateWorkflowServices(server *grpc.Server, apiServer *api.Server) {
	pb.RegisterGoPayAppWorkflowServiceServer(server, api.NewGoPayAppWorkflowService(apiServer))
}
