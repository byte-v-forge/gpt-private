//go:build private_plugins

package activities

import "orchestrator/pb"

func paymentPrepareData(resp *pb.PrepareGoPayResponse) map[string]any {
	if resp == nil {
		return map[string]any{"response_present": false}
	}
	return map[string]any{
		"response_present":         true,
		"success":                  resp.GetSuccess(),
		"error_message":            resp.GetErrorMessage(),
		"flow_id":                  resp.GetFlowId(),
		"snap_token_present":       resp.GetSnapToken() != "",
		"checkout_url":             resp.GetCheckoutUrl(),
		"checkout_session_id":      resp.GetCheckoutSessionId(),
		"retryable_fresh_checkout": resp.GetRetryableFreshCheckout(),
		"checkout_attempt":         resp.GetCheckoutAttempt(),
		"stage":                    resp.GetStage(),
	}
}
