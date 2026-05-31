//go:build private_plugins

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/byte-v-forge/gpt/pkg/gptplugin"
)

type rawN8NActionRequest struct {
	JobID             string         `json:"job_id"`
	AccountID         string         `json:"account_id"`
	N8NExecutionID    string         `json:"n8n_execution_id"`
	FlowID            string         `json:"flow_id"`
	SnapToken         string         `json:"snap_token"`
	CheckoutURL       string         `json:"checkout_url"`
	CheckoutSessionID string         `json:"checkout_session_id"`
	StateJSON         string         `json:"state_json"`
	UseAccountToken   bool           `json:"use_account_token"`
	OTPRequired       bool           `json:"otp_required"`
	OTPSource         string         `json:"otp_source"`
	OTP               string         `json:"otp"`
	Channel           string         `json:"channel"`
	Target            string         `json:"target"`
	ResumeURL         string         `json:"resume_url"`
	OTPIssuedAfter    int64          `json:"otp_issued_after_unix"`
	OTPReceivedAt     int64          `json:"otp_received_at_unix"`
	OTPTimeoutSeconds int32          `json:"otp_timeout_seconds"`
	OTPRetryAttempt   int32          `json:"otp_retry_attempt"`
	GopayAccountID    string         `json:"gopay_account_id"`
	Operation         string         `json:"operation"`
	Phone             string         `json:"phone"`
	OTPChannel        string         `json:"otp_channel"`
	ActivationID      string         `json:"activation_id"`
	WAPhone           string         `json:"wa_phone"`
	ChargeRef         string         `json:"charge_ref"`
	PlusTrialEligible bool           `json:"plus_trial_eligible"`
	PlusTrialChecked  bool           `json:"plus_trial_checked"`
	PlusActive        bool           `json:"plus_active"`
	FailureCount      int32          `json:"failure_count"`
	ProxyHash         string         `json:"proxy_hash"`
	ProxyURL          string         `json:"proxy_url"`
	DeviceFingerprint string         `json:"device_fingerprint"`
	Reason            string         `json:"reason"`
	ErrorMessage      string         `json:"error_message"`
	Data              map[string]any `json:"data"`
}

func (s *Server) InvokeN8NAction(ctx context.Context, call gptplugin.N8NActionCall) (any, error) {
	actionID := strings.ToUpper(strings.TrimSpace(call.ActionID))
	subPath := strings.Trim(strings.TrimSpace(call.SubPath), "/")
	var req rawN8NActionRequest
	if len(call.RawJSON) > 0 {
		if err := json.Unmarshal(call.RawJSON, &req); err != nil {
			return nil, fmt.Errorf("decode n8n action request: %w", err)
		}
	}
	switch actionID {
	case actionGoPayPayment:
		return s.invokeN8NGoPayPayment(ctx, subPath, req)
	case actionGoPayQRISPaymentActivate:
		return s.invokeN8NGoPayQRISPayment(ctx, subPath, req)
	default:
		return nil, fmt.Errorf("unsupported raw n8n action: %s", actionID)
	}
}

func (s *Server) invokeN8NGoPayPayment(ctx context.Context, action string, req rawN8NActionRequest) (any, error) {
	switch action {
	case "proxy-settings":
		return s.N8NGoPayDynamicProxySettings(ctx, req.JobID, req.AccountID, req.N8NExecutionID, actionGoPayPayment)
	case "record-proxy":
		return s.RecordN8NGoPayDynamicProxy(ctx, req.JobID, req.AccountID, req.N8NExecutionID, req.ProxyURL, req.Data, actionGoPayPayment)
	case "fail-proxy":
		return s.FailN8NGoPayDynamicProxy(ctx, req.JobID, req.AccountID, req.N8NExecutionID, req.ErrorMessage, req.Data)
	case "resolve-account":
		return s.ResolveN8NGoPayPaymentAccount(ctx, req.JobID, req.AccountID, req.N8NExecutionID)
	case "probe-plus-trial":
		return s.ProbeN8NGoPayPaymentPlusTrial(ctx, req.JobID, req.AccountID, req.N8NExecutionID)
	case "prepare-checkout":
		return s.PrepareN8NGoPayPaymentCheckout(ctx, req.JobID, req.AccountID, req.N8NExecutionID)
	case "prepare-link":
		return s.PrepareN8NGoPayPaymentLink(ctx, req.JobID, req.AccountID, req.N8NExecutionID, req.FlowID, req.CheckoutURL, req.CheckoutSessionID, req.StateJSON)
	case "check-channel-otp":
		return s.CheckN8NGoPayPaymentChannelOTP(ctx, req.JobID, req.AccountID, req.N8NExecutionID, req.ActivationID, req.OTPIssuedAfter)
	case "finish":
		return s.FinishN8NGoPayPayment(ctx, req.JobID, req.AccountID, req.N8NExecutionID, req.ChargeRef, req.Data)
	case "fail":
		return s.FailN8NGoPayPayment(ctx, req.JobID, req.N8NExecutionID, req.FlowID, req.ErrorMessage, req.Data)
	default:
		return nil, fmt.Errorf("unsupported gopay payment action: %s", action)
	}
}

func (s *Server) invokeN8NGoPayQRISPayment(ctx context.Context, action string, req rawN8NActionRequest) (any, error) {
	switch action {
	case "proxy-settings":
		return s.N8NGoPayDynamicProxySettings(ctx, req.JobID, req.AccountID, req.N8NExecutionID, actionGoPayQRISPaymentActivate)
	case "record-proxy":
		return s.RecordN8NGoPayDynamicProxy(ctx, req.JobID, req.AccountID, req.N8NExecutionID, req.ProxyURL, req.Data, actionGoPayQRISPaymentActivate)
	case "fail-proxy":
		return s.FailN8NGoPayDynamicProxy(ctx, req.JobID, req.AccountID, req.N8NExecutionID, req.ErrorMessage, req.Data)
	case "resolve-account":
		return s.ResolveN8NGoPayQRISPaymentAccount(ctx, req.JobID, req.AccountID, req.N8NExecutionID)
	case "probe-plus-trial":
		return s.ProbeN8NGoPayQRISPaymentPlusTrial(ctx, req.JobID, req.AccountID, req.N8NExecutionID)
	case "prepare-checkout":
		return s.PrepareN8NGoPayQRISPaymentCheckout(ctx, req.JobID, req.AccountID, req.N8NExecutionID)
	case "prepare-link":
		return s.PrepareN8NGoPayQRISPaymentLink(ctx, req.JobID, req.AccountID, req.N8NExecutionID, req.FlowID, req.CheckoutURL, req.CheckoutSessionID, req.StateJSON)
	case "check-manual-payment":
		return s.CheckN8NGoPayQRISManualPayment(ctx, req.JobID, req.AccountID, req.N8NExecutionID, req.FlowID)
	case "finish":
		return s.FinishN8NGoPayQRISPaymentActivate(ctx, req.JobID, req.AccountID, req.N8NExecutionID, req.ChargeRef, req.PlusTrialEligible, req.PlusTrialChecked, req.PlusActive, req.Data)
	case "fail":
		return s.FailN8NGoPayQRISPaymentActivate(ctx, req.JobID, req.N8NExecutionID, req.FlowID, req.ErrorMessage, req.Data)
	default:
		return nil, fmt.Errorf("unsupported gopay qris action: %s", action)
	}
}
