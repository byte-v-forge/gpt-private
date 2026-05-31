//go:build private_plugins

package activities

import "orchestrator/pb"

var privateGoPayClient pb.GopayAppServiceClient

func ConfigureGoPayDependencies(client pb.GopayAppServiceClient) {
	privateGoPayClient = client
}

func (s *Server) goPayClient() pb.GopayAppServiceClient {
	return privateGoPayClient
}
