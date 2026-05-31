//go:build private_plugins

package api

import (
	"context"

	"orchestrator/pb"
)

var privateGoPayClient pb.GopayAppServiceClient

func ConfigureGoPayClient(client pb.GopayAppServiceClient) {
	privateGoPayClient = client
}

func (s *Server) goPayClient() pb.GopayAppServiceClient {
	return privateGoPayClient
}

type goPayAppWorkflowService struct {
	pb.UnimplementedGoPayAppWorkflowServiceServer
	server *Server
}

func NewGoPayAppWorkflowService(server *Server) pb.GoPayAppWorkflowServiceServer {
	return &goPayAppWorkflowService{server: server}
}

func (s *goPayAppWorkflowService) RunGoPayApp(ctx context.Context, req *pb.GoPayAppRequest) (*pb.GoPayAppResponse, error) {
	return s.server.RunGoPayApp(ctx, req)
}

func (s *goPayAppWorkflowService) RunGoPayPayment(ctx context.Context, req *pb.GoPayPaymentRequest) (*pb.GoPayPaymentResponse, error) {
	return s.server.RunGoPayPayment(ctx, req)
}

func (s *goPayAppWorkflowService) RunGoPayQRISPaymentActivate(ctx context.Context, req *pb.GoPayQRISPaymentActivateRequest) (*pb.GoPayPaymentResponse, error) {
	return s.server.RunGoPayQRISPaymentActivate(ctx, req)
}

func (s *goPayAppWorkflowService) RunGoPayWAPayment(ctx context.Context, req *pb.GoPayWAPaymentRequest) (*pb.GoPayPaymentResponse, error) {
	return s.server.RunGoPayWAPayment(ctx, req)
}

func (s *goPayAppWorkflowService) RetryGoPayPaymentRebind(ctx context.Context, req *pb.GoPayPaymentRebindRequest) (*pb.GoPayPaymentResponse, error) {
	return s.server.RetryGoPayPaymentRebind(ctx, req)
}

func (s *goPayAppWorkflowService) ConfirmManualGoPayPayment(ctx context.Context, req *pb.ConfirmManualGoPayPaymentRequest) (*pb.ConfirmManualGoPayPaymentResponse, error) {
	return s.server.ConfirmManualGoPayPayment(ctx, req)
}

func (s *goPayAppWorkflowService) GetGopayAccount(ctx context.Context, req *pb.GetGopayAccountRequest) (*pb.GetGopayAccountResponse, error) {
	if s.server.goPayClient() == nil {
		return &pb.GetGopayAccountResponse{ErrorMessage: "gopay-app client not configured"}, nil
	}
	return s.server.goPayClient().GetGopayAccount(ctx, req)
}

func (s *goPayAppWorkflowService) ListGopayAccounts(ctx context.Context, req *pb.ListGopayAccountsRequest) (*pb.ListGopayAccountsResponse, error) {
	if s.server.goPayClient() == nil {
		return &pb.ListGopayAccountsResponse{ErrorMessage: "gopay-app client not configured"}, nil
	}
	return s.server.goPayClient().ListGopayAccounts(ctx, req)
}

func (s *goPayAppWorkflowService) DeleteGopayAccount(ctx context.Context, req *pb.DeleteGopayAccountRequest) (*pb.DeleteGopayAccountResponse, error) {
	if s.server.goPayClient() == nil {
		return &pb.DeleteGopayAccountResponse{ErrorMessage: "gopay-app client not configured"}, nil
	}
	return s.server.goPayClient().DeleteGopayAccount(ctx, req)
}
