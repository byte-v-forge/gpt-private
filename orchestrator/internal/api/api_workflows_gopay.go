//go:build private_plugins

package api

import (
	"context"
	"strings"

	"orchestrator/internal/jobstatus"
	"orchestrator/pb"
)

func (s *Server) RunGoPayApp(context.Context, *pb.GoPayAppRequest) (*pb.GoPayAppResponse, error) {
	return &pb.GoPayAppResponse{ErrorMessage: "gopay-app is n8n-only; use GPT dashboard BFF workflow endpoint"}, nil
}

func (s *Server) RunGoPayPayment(context.Context, *pb.GoPayPaymentRequest) (*pb.GoPayPaymentResponse, error) {
	return &pb.GoPayPaymentResponse{ErrorMessage: "gopay-payment is n8n-only; use GPT dashboard BFF workflow endpoint"}, nil
}

func (s *Server) RunGoPayQRISPaymentActivate(context.Context, *pb.GoPayQRISPaymentActivateRequest) (*pb.GoPayPaymentResponse, error) {
	return &pb.GoPayPaymentResponse{ErrorMessage: "gopay-qris-payment-activate is n8n-only; use GPT dashboard BFF workflow endpoint"}, nil
}

func (s *Server) RunGoPayWAPayment(context.Context, *pb.GoPayWAPaymentRequest) (*pb.GoPayPaymentResponse, error) {
	return &pb.GoPayPaymentResponse{ErrorMessage: "gopay-wa-payment is n8n-only; use GPT dashboard BFF workflow endpoint"}, nil
}

func (s *Server) RetryGoPayPaymentRebind(context.Context, *pb.GoPayPaymentRebindRequest) (*pb.GoPayPaymentResponse, error) {
	return &pb.GoPayPaymentResponse{ErrorMessage: "gopay-payment-rebind is n8n-only; use GPT dashboard BFF workflow endpoint"}, nil
}

func (s *Server) ConfirmManualGoPayPayment(ctx context.Context, req *pb.ConfirmManualGoPayPaymentRequest) (*pb.ConfirmManualGoPayPaymentResponse, error) {
	jobID := strings.TrimSpace(req.GetJobId())
	if jobID == "" {
		return &pb.ConfirmManualGoPayPaymentResponse{Success: false, ErrorMessage: "job_id is required"}, nil
	}
	job, err := s.getJob(ctx, jobID)
	if err != nil {
		return &pb.ConfirmManualGoPayPaymentResponse{Success: false, JobId: jobID, ErrorMessage: err.Error()}, nil
	}
	if job.Status != jobstatus.Running {
		return &pb.ConfirmManualGoPayPaymentResponse{Success: false, JobId: jobID, ErrorMessage: "job is not running: " + job.Status}, nil
	}
	if job.Action != actionGoPayQRISPaymentActivate && job.Action != actionGoPayPayment {
		return &pb.ConfirmManualGoPayPaymentResponse{Success: false, JobId: jobID, ErrorMessage: "job does not accept manual gopay payment confirmation: " + job.Action}, nil
	}
	if job.LastStep != stepGoPayPayment {
		return &pb.ConfirmManualGoPayPaymentResponse{Success: false, JobId: jobID, ErrorMessage: "job is not waiting for gopay payment confirmation: " + job.LastStep}, nil
	}
	if err := s.setJobParams(ctx, jobID, map[string]string{manualGoPayPaymentConfirmParam: "true"}); err != nil {
		return &pb.ConfirmManualGoPayPaymentResponse{Success: false, JobId: jobID, ErrorMessage: err.Error()}, nil
	}
	return &pb.ConfirmManualGoPayPaymentResponse{Success: true, JobId: jobID}, nil
}

func goPayQRISPaymentJobParams(req *pb.GoPayQRISPaymentActivateRequest) map[string]string {
	params := map[string]string{
		"activation_mode":       "qris_payment",
		"payment_type":          "qris",
		"tokenization":          "qris",
		"otp_channel":           "not_required",
		"uses_wa":               "false",
		"uses_gopay_app_flow":   "false",
		"manual_confirmation":   "true",
		"manual_payment_button": "true",
	}
	if value := strings.TrimSpace(req.GetAccountId()); value != "" {
		params["account_id"] = value
	}
	if value := strings.TrimSpace(req.GetSourceJobId()); value != "" {
		params["source_job_id"] = value
	}
	if value := goPayAppAccountID(req.GetGopayAccountId()); value != "" {
		params["gopay_account_id"] = value
	}
	return params
}

func goPayPaymentJobParams(req *pb.GoPayPaymentRequest) map[string]string {
	params := map[string]string{
		"tokenization": firstNonEmpty(strings.TrimSpace(req.GetTokenization()), "true"),
	}
	if value := strings.TrimSpace(req.GetAccountId()); value != "" {
		params["account_id"] = value
	}
	if value := strings.TrimSpace(req.GetSourceJobId()); value != "" {
		params["source_job_id"] = value
	}
	if value := goPayAppAccountID(req.GetGopayAccountId()); value != "" {
		params["gopay_account_id"] = value
	}
	return params
}
