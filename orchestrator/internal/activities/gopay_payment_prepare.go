//go:build private_plugins

package activities

import (
	"context"
	"fmt"
	"strings"

	"orchestrator/internal/gptaccount"
	"orchestrator/pb"
)

func (s *Server) prepareGoPayPayment(ctx context.Context, step activityStep, input GoPayActivityInput, account *pb.Account) (output GoPayPaymentPrepareOutput, err error) {
	accountID := gptaccount.ID(account)
	sessionToken := strings.TrimSpace(input.GetSessionToken())
	if sessionToken == "" {
		sessionToken = s.cachedChatGPTSessionToken(ctx, accountID)
	}
	accessToken := strings.TrimSpace(input.GetAccessToken())
	if accessToken == "" {
		accessToken = s.cachedChatGPTAccessToken(ctx, accountID)
	}
	tokenization := strings.TrimSpace(input.GetTokenization())
	suppliedCheckoutURL := strings.TrimSpace(input.GetCheckoutUrl())
	suppliedCheckoutSessionID := strings.TrimSpace(input.GetCheckoutSessionId())
	checkoutURL := ""
	checkoutSessionID := ""
	gopayPhone := normalizeIndonesiaPhone(input.GetGopayPhone())
	stateJSON := normalizeGoPayWorkflowStateJSON(input.GetStateJson())

	data := map[string]any{
		"account_id":             accountID,
		"session_token_present":  sessionToken != "",
		"proxy_present":          strings.TrimSpace(input.GetProxyUrl()) != "",
		"access_token_present":   accessToken != "",
		"tokenization":           tokenization,
		"checkout_url_present":   suppliedCheckoutURL != "",
		"checkout_session_id":    suppliedCheckoutSessionID,
		"checkout_reuse_blocked": suppliedCheckoutURL != "" || suppliedCheckoutSessionID != "",
		"gopay_phone_present":    gopayPhone != "",
		"prepared_flow_present":  false,
		"payment_prepare_called": false,
	}
	output = GoPayPaymentPrepareOutput{
		UseAccountToken: false,
		StateJson:       stateJSON,
	}
	defer func() {
		output.StateJson = stateJSON
		output.Data = protoData(data)
	}()
	if sessionToken == "" && accessToken == "" {
		return output, fmt.Errorf("session_token or access_token is required")
	}

	step.progress("preparing gopay payment before user link", protoData(map[string]any{
		"tokenization":              tokenization,
		"supplied_checkout_present": suppliedCheckoutURL != "" || suppliedCheckoutSessionID != "",
		"checkout_reuse_blocked":    suppliedCheckoutURL != "" || suppliedCheckoutSessionID != "",
		"gopay_phone_present":       gopayPhone != "",
	}))
	credential, err := s.paymentCredentialWithProxy(ctx, accountID, sessionToken, accessToken, input.GetProxyUrl())
	if err != nil {
		return output, err
	}
	prepared, err := s.paymentClient.PrepareGoPay(ctx, &pb.PrepareGoPayRequest{
		Credential:        credential,
		Tokenization:      tokenization,
		CheckoutUrl:       checkoutURL,
		CheckoutSessionId: checkoutSessionID,
		GopayPhone:        gopayPhone,
		GopayCountryCode:  input.GetCountryCode(),
	})
	data["payment_prepare_called"] = true
	data["payment_prepare"] = paymentPrepareData(prepared)
	if prepared != nil {
		applyGoPayPaymentPrepareResponse(&output, prepared)
		data["prepared_flow_present"] = output.GetFlowId() != ""
	}
	step.progress("gopay payment prepared", protoData(map[string]any{
		"success":            prepared != nil && prepared.GetSuccess(),
		"flow_id_present":    output.GetFlowId() != "",
		"snap_token_present": output.GetSnapToken() != "",
	}))
	if err != nil {
		return output, err
	}
	if prepared == nil {
		return output, fmt.Errorf("payment prepare returned empty response")
	}
	if !prepared.GetSuccess() {
		return output, fmt.Errorf("payment prepare failed: %s", prepared.GetErrorMessage())
	}
	if output.GetFlowId() == "" {
		return output, fmt.Errorf("payment prepare returned empty flow_id")
	}
	return output, nil
}

func applyGoPayPaymentPrepareResponse(output *GoPayPaymentPrepareOutput, resp *pb.PrepareGoPayResponse) {
	if output == nil || resp == nil {
		return
	}
	output.FlowId = resp.GetFlowId()
	output.SnapToken = resp.GetSnapToken()
	output.CheckoutUrl = resp.GetCheckoutUrl()
	output.CheckoutSessionId = resp.GetCheckoutSessionId()
	output.RetryableFreshCheckout = resp.GetRetryableFreshCheckout()
	output.CheckoutAttempt = resp.GetCheckoutAttempt()
	output.Stage = resp.GetStage()
}
