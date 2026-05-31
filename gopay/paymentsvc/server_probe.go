package paymentsvc

import (
	"context"

	"github.com/byte-v-forge/common-lib/stringx"
	"github.com/byte-v-forge/gpt-private/gopay/pb"
)

func (s *Server) ProbeTier(ctx context.Context, req *pb.ProbeTierPaymentRequest) (*pb.ProbeTierPaymentResponse, error) {
	cred := requestCredential(req.GetCredential())
	if cred.empty() {
		return &pb.ProbeTierPaymentResponse{Success: false, ErrorMessage: "credential is required"}, nil
	}
	result, err := s.probeTier(ctx, cred)
	if err != nil {
		return &pb.ProbeTierPaymentResponse{Success: false, ErrorMessage: truncateError(err)}, nil
	}
	return &pb.ProbeTierPaymentResponse{
		Success:      result.ErrorMessage == "",
		ErrorMessage: result.ErrorMessage,
		Checked:      result.Checked,
		Tier:         stringx.FirstNonEmpty(result.Tier, result.PlanType),
		PlusActive:   result.PlusActive,
		Source:       stringx.FirstNonEmpty(result.Source, "auth_session"),
	}, nil
}

func (s *Server) ProbePlusTrial(ctx context.Context, req *pb.ProbePlusTrialPaymentRequest) (*pb.ProbePlusTrialPaymentResponse, error) {
	cred := requestCredential(req.GetCredential())
	if cred.empty() {
		return &pb.ProbePlusTrialPaymentResponse{Success: false, ErrorMessage: "session_token or access_token is required"}, nil
	}
	result, err := s.probePlusTrial(ctx, cred)
	if err != nil {
		return &pb.ProbePlusTrialPaymentResponse{Success: false, ErrorMessage: truncateError(err)}, nil
	}
	return &pb.ProbePlusTrialPaymentResponse{
		Success:           result.ErrorMessage == "",
		ErrorMessage:      result.ErrorMessage,
		Checked:           result.Checked,
		PlusTrialEligible: result.PlusTrialEligible,
		Amount:            result.Amount,
		Currency:          result.Currency,
		Source:            result.Source,
		CheckoutUrl:       result.CheckoutURL,
		CheckoutSessionId: result.CheckoutSessionID,
		PlusActive:        result.PlusActive,
		PlanType:          result.PlanType,
	}, nil
}

func (s *Server) CreateCheckoutLink(ctx context.Context, req *pb.CreateCheckoutLinkRequest) (*pb.CreateCheckoutLinkResponse, error) {
	cred := requestCredential(req.GetCredential())
	if cred.empty() {
		return &pb.CreateCheckoutLinkResponse{Success: false, ErrorMessage: "session_token or access_token is required"}, nil
	}
	ch, err := s.newCharger(ctx, cred, "", "", "", defaultTokenization)
	if err != nil {
		return &pb.CreateCheckoutLinkResponse{Success: false, ErrorMessage: truncateError(err)}, nil
	}
	defer ch.close()
	csID, err := ch.createCheckout(ctx)
	if err != nil {
		return &pb.CreateCheckoutLinkResponse{Success: false, ErrorMessage: truncateError(err)}, nil
	}
	return &pb.CreateCheckoutLinkResponse{Success: true, CheckoutUrl: ch.checkoutURL, CheckoutSessionId: csID}, nil
}
