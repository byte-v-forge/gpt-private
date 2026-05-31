package paymentsvc

import (
	"context"
	"strings"

	"github.com/byte-v-forge/gpt-private/gopay/pb"
)

func (s *Server) PrepareGoPay(ctx context.Context, req *pb.PrepareGoPayRequest) (*pb.PrepareGoPayResponse, error) {
	cred := requestCredential(req.GetCredential())
	if cred.empty() {
		return &pb.PrepareGoPayResponse{Success: false, ErrorMessage: "session_token or access_token is required"}, nil
	}
	ch, err := s.newCharger(ctx, cred, req.GetGopayPhone(), req.GetGopayCountryCode(), "", req.GetTokenization())
	if err != nil {
		return &pb.PrepareGoPayResponse{Success: false, ErrorMessage: truncateError(err)}, nil
	}
	state, err := ch.prepareUntilLinking(ctx, req.GetCheckoutSessionId(), req.GetCheckoutUrl())
	if err != nil {
		ch.close()
		return &pb.PrepareGoPayResponse{Success: false, ErrorMessage: truncateError(err)}, nil
	}
	flowID := s.flows.put(&pendingFlow{charger: ch, state: state})
	return &pb.PrepareGoPayResponse{
		Success:           true,
		FlowId:            flowID,
		SnapToken:         stringAt(state, "snap_token"),
		CheckoutUrl:       stringAt(state, "checkout_url"),
		CheckoutSessionId: stringAt(state, "cs_id"),
		CheckoutAttempt:   int32(intAt(state, "checkout_attempt")),
		Stage:             stringAt(state, "state"),
	}, nil
}

func (s *Server) PrepareGoPayCheckout(ctx context.Context, req *pb.PrepareGoPayCheckoutRequest) (*pb.PrepareGoPayResponse, error) {
	cred := requestCredential(req.GetCredential())
	if cred.empty() {
		return &pb.PrepareGoPayResponse{Success: false, ErrorMessage: "session_token or access_token is required", Stage: "checkout"}, nil
	}
	ch, err := s.newCharger(ctx, cred, req.GetGopayPhone(), req.GetGopayCountryCode(), "", req.GetTokenization())
	if err != nil {
		return &pb.PrepareGoPayResponse{Success: false, ErrorMessage: truncateError(err), Stage: "checkout"}, nil
	}
	csID, state, err := ch.prepareCheckout(ctx, req.GetCheckoutSessionId(), req.GetCheckoutUrl(), 1)
	if err != nil {
		ch.close()
		return &pb.PrepareGoPayResponse{Success: false, ErrorMessage: truncateError(err), Stage: "checkout"}, nil
	}
	flowID := s.flows.put(&pendingFlow{charger: ch, state: state})
	return &pb.PrepareGoPayResponse{
		Success:           true,
		FlowId:            flowID,
		CheckoutUrl:       stringAt(state, "checkout_url"),
		CheckoutSessionId: csID,
		CheckoutAttempt:   int32(intAt(state, "checkout_attempt")),
		Stage:             "checkout",
	}, nil
}

func (s *Server) RefreshPrepareGoPayCheckout(ctx context.Context, req *pb.RefreshPrepareGoPayCheckoutRequest) (*pb.PrepareGoPayResponse, error) {
	flowID := strings.TrimSpace(req.GetFlowId())
	if flowID == "" {
		return &pb.PrepareGoPayResponse{Success: false, ErrorMessage: "flow_id is required", Stage: "checkout_refresh"}, nil
	}
	flow := s.flows.get(flowID)
	if flow == nil {
		return &pb.PrepareGoPayResponse{Success: false, ErrorMessage: "prepared payment flow not found", Stage: "checkout_refresh"}, nil
	}
	nextAttempt := int(intAt(flow.state, "checkout_attempt")) + 1
	if nextAttempt <= 1 {
		nextAttempt = 2
	}
	csID, state, err := flow.charger.prepareCheckout(ctx, "", "", nextAttempt)
	if err != nil {
		if failed := s.flows.pop(flowID); failed != nil {
			failed.close()
		}
		return &pb.PrepareGoPayResponse{Success: false, ErrorMessage: truncateError(err), FlowId: flowID, Stage: "checkout_refresh"}, nil
	}
	flow.state = state
	return &pb.PrepareGoPayResponse{
		Success:           true,
		FlowId:            flowID,
		CheckoutUrl:       stringAt(state, "checkout_url"),
		CheckoutSessionId: csID,
		CheckoutAttempt:   int32(intAt(state, "checkout_attempt")),
		Stage:             "checkout_refresh",
	}, nil
}

func (s *Server) PrepareGoPayLink(ctx context.Context, req *pb.PrepareGoPayLinkRequest) (*pb.PrepareGoPayResponse, error) {
	flowID := strings.TrimSpace(req.GetFlowId())
	if flowID == "" {
		return &pb.PrepareGoPayResponse{Success: false, ErrorMessage: "flow_id is required", Stage: "link"}, nil
	}
	flow := s.flows.get(flowID)
	if flow == nil {
		return &pb.PrepareGoPayResponse{Success: false, ErrorMessage: "prepared payment flow not found", FlowId: flowID, Stage: "link"}, nil
	}
	csID := stringAt(flow.state, "cs_id")
	if csID == "" {
		return &pb.PrepareGoPayResponse{Success: false, ErrorMessage: "prepared payment flow is missing checkout_session_id", FlowId: flowID, Stage: "link"}, nil
	}
	state, err := flow.charger.prepareCheckoutSessionUntilLinking(ctx, csID)
	if err != nil {
		if isChatGPTApproveBlocked(err) {
			return &pb.PrepareGoPayResponse{
				Success:                false,
				ErrorMessage:           truncateError(err),
				FlowId:                 flowID,
				CheckoutUrl:            stringAt(flow.state, "checkout_url"),
				CheckoutSessionId:      csID,
				RetryableFreshCheckout: true,
				CheckoutAttempt:        int32(intAt(flow.state, "checkout_attempt")),
				Stage:                  "link",
			}, nil
		}
		if failed := s.flows.pop(flowID); failed != nil {
			failed.close()
		}
		return &pb.PrepareGoPayResponse{
			Success:           false,
			ErrorMessage:      truncateError(err),
			FlowId:            flowID,
			CheckoutUrl:       stringAt(flow.state, "checkout_url"),
			CheckoutSessionId: csID,
			CheckoutAttempt:   int32(intAt(flow.state, "checkout_attempt")),
			Stage:             "link",
		}, nil
	}
	state["checkout_attempt"] = intAt(flow.state, "checkout_attempt")
	flow.state = state
	return &pb.PrepareGoPayResponse{
		Success:           true,
		FlowId:            flowID,
		SnapToken:         stringAt(state, "snap_token"),
		CheckoutUrl:       stringAt(state, "checkout_url"),
		CheckoutSessionId: stringAt(state, "cs_id"),
		CheckoutAttempt:   int32(intAt(state, "checkout_attempt")),
		Stage:             "link",
	}, nil
}
