//go:build private_plugins

package api

import (
	"context"
	"fmt"
	"strings"

	"orchestrator/internal/jobstatus"
	"orchestrator/pb"
)

type n8nGoPayPaymentStepResult struct {
	JobID              string         `json:"job_id"`
	AccountID          string         `json:"account_id,omitempty"`
	N8NExecutionID     string         `json:"n8n_execution_id,omitempty"`
	Action             string         `json:"action"`
	Step               string         `json:"step"`
	Success            bool           `json:"success"`
	Checked            bool           `json:"checked,omitempty"`
	PlusTrialEligible  bool           `json:"plus_trial_eligible,omitempty"`
	PlusTrialChecked   bool           `json:"plus_trial_checked,omitempty"`
	PlusActive         bool           `json:"plus_active,omitempty"`
	GopayAccountID     string         `json:"gopay_account_id,omitempty"`
	Phone              string         `json:"phone,omitempty"`
	ActivationID       string         `json:"activation_id,omitempty"`
	StateJSON          string         `json:"state_json,omitempty"`
	FlowID             string         `json:"flow_id,omitempty"`
	CheckoutURL        string         `json:"checkout_url,omitempty"`
	CheckoutSessionID  string         `json:"checkout_session_id,omitempty"`
	UseAccountToken    bool           `json:"use_account_token,omitempty"`
	Ready              bool           `json:"ready,omitempty"`
	AccountTokenReady  bool           `json:"account_token_ready,omitempty"`
	SignupComplete     bool           `json:"signup_complete,omitempty"`
	PhoneAccepted      bool           `json:"phone_accepted,omitempty"`
	DeviceProxyMatched bool           `json:"device_proxy_matched,omitempty"`
	RetryableFailure   bool           `json:"retryable_failure,omitempty"`
	RotatableFailure   bool           `json:"rotatable_failure,omitempty"`
	FailureCount       int32          `json:"failure_count,omitempty"`
	MaxFailures        int32          `json:"max_failures,omitempty"`
	ProxyHash          string         `json:"proxy_hash,omitempty"`
	DeviceFingerprint  string         `json:"device_fingerprint,omitempty"`
	OTPRequired        bool           `json:"otp_required,omitempty"`
	OTPChannel         string         `json:"otp_channel,omitempty"`
	OTPFound           bool           `json:"otp_found,omitempty"`
	OTPSource          string         `json:"otp_source,omitempty"`
	OTPIssuedAfterUnix int64          `json:"otp_issued_after_unix,omitempty"`
	OTPTimeoutSeconds  int32          `json:"otp_timeout_seconds,omitempty"`
	ChargeRef          string         `json:"charge_ref,omitempty"`
	SnapToken          string         `json:"snap_token,omitempty"`
	SnapTokenPresent   bool           `json:"snap_token_present,omitempty"`
	Data               map[string]any `json:"data,omitempty"`
}

func (s *Server) ResolveN8NGoPayPaymentAccount(ctx context.Context, jobID string, accountID string, n8nExecutionID string) (any, error) {
	jobID, accountID, n8nExecutionID = normalizeN8NGoPayQRISIDs(jobID, accountID, n8nExecutionID)
	if err := s.bindN8NGoPayExecution(ctx, jobID, n8nExecutionID); err != nil {
		return nil, err
	}
	params, err := s.jobStore.Params(ctx, jobID)
	if err != nil {
		return nil, err
	}
	data := goPayPaymentBaseData(params)
	gopayAccountID := goPayAppAccountID(params["gopay_account_id"])
	if gopayAccountID == "" {
		err := fmt.Errorf("gopay_account_id is required")
		return &n8nGoPayPaymentStepResult{JobID: jobID, AccountID: accountID, N8NExecutionID: n8nExecutionID, Action: actionGoPayPayment, Step: "resolve_account", Success: false, Data: data}, s.markActionFailed(ctx, jobID, "resolve_account", jobstatus.FailedFinal, false, false, err, data)
	}
	account, err := s.activities.ResolveAccountFromJobActivity(ctx, pb.ResolveAccountInput{AccountId: firstNonEmpty(accountID, params["account_id"]), SourceJobId: params["source_job_id"]})
	resolvedAccountID := account.GetAccountId()
	result := &n8nGoPayPaymentStepResult{JobID: jobID, AccountID: resolvedAccountID, N8NExecutionID: n8nExecutionID, Action: actionGoPayPayment, Step: "resolve_account", Success: err == nil, GopayAccountID: gopayAccountID, Data: data}
	if err != nil {
		return result, s.markActionFailed(ctx, jobID, "resolve_account", jobstatus.FailedRetryable, false, true, err, data)
	}
	if err := s.jobStore.SetAccountID(ctx, jobID, resolvedAccountID); err != nil {
		return result, s.markActionFailed(ctx, jobID, "resolve_account", jobstatus.FailedRetryable, false, true, err, data)
	}
	data["account_id"] = resolvedAccountID
	return result, nil
}

func (s *Server) ProbeN8NGoPayPaymentPlusTrial(ctx context.Context, jobID string, accountID string, n8nExecutionID string) (any, error) {
	jobID, accountID, n8nExecutionID = normalizeN8NGoPayQRISIDs(jobID, accountID, n8nExecutionID)
	if err := s.bindN8NGoPayExecution(ctx, jobID, n8nExecutionID); err != nil {
		return nil, err
	}
	probe, err := s.activities.ProbePlusTrialAtomicActivity(ctx, pb.ProbePlusTrialActivityInput{JobId: jobID, AccountId: accountID, ProxyUrl: s.protocolProxyURL(ctx, jobID)})
	data := structMap(probe.GetData())
	result := &n8nGoPayPaymentStepResult{JobID: jobID, AccountID: accountID, N8NExecutionID: n8nExecutionID, Action: actionGoPayPayment, Step: stepProbePlusTrial, Success: err == nil, Checked: probe.GetChecked(), PlusTrialEligible: probe.GetPlusTrialEligible(), PlusTrialChecked: probe.GetChecked(), PlusActive: probe.GetPlusActive(), Data: data}
	if err != nil {
		return result, s.markActionFailed(ctx, jobID, stepProbePlusTrial, jobstatus.FailedRetryable, false, true, err, data)
	}
	return result, nil
}

func (s *Server) PrepareN8NGoPayPaymentCheckout(ctx context.Context, jobID string, accountID string, n8nExecutionID string) (any, error) {
	return s.prepareN8NGoPayPayment(ctx, stepGoPayPaymentPrepareCheckout, jobID, accountID, n8nExecutionID, "", "", "", "{}")
}

func (s *Server) PrepareN8NGoPayPaymentLink(ctx context.Context, jobID string, accountID string, n8nExecutionID string, flowID string, checkoutURL string, checkoutSessionID string, stateJSON string) (any, error) {
	return s.prepareN8NGoPayPayment(ctx, stepGoPayPaymentPrepareLink, jobID, accountID, n8nExecutionID, flowID, checkoutURL, checkoutSessionID, stateJSON)
}

func (s *Server) CheckN8NGoPayPaymentChannelOTP(ctx context.Context, jobID string, accountID string, n8nExecutionID string, activationID string, issuedAfterUnix int64) (any, error) {

	jobID, accountID, n8nExecutionID = normalizeN8NGoPayQRISIDs(jobID, accountID, n8nExecutionID)
	if err := s.bindN8NGoPayExecution(ctx, jobID, n8nExecutionID); err != nil {
		return nil, err
	}
	data := goPayOTPCheckData("sms", strings.TrimSpace(activationID), issuedAfterUnix)
	return goPayPaymentChannelOTPCheckResult(jobID, accountID, n8nExecutionID, activationID, issuedAfterUnix, "", false, data), nil
}

func (s *Server) FinishN8NGoPayPayment(ctx context.Context, jobID string, accountID string, n8nExecutionID string, chargeRef string, data map[string]any) (any, error) {
	jobID, accountID, n8nExecutionID = normalizeN8NGoPayQRISIDs(jobID, accountID, n8nExecutionID)
	if err := s.bindN8NGoPayExecution(ctx, jobID, n8nExecutionID); err != nil {
		return nil, err
	}
	result := goPayPaymentBaseData(nil)
	for key, value := range data {
		result[key] = value
	}
	result["account_id"] = accountID
	result["payment_completed"] = strings.TrimSpace(chargeRef) != ""
	result["charge_ref"] = strings.TrimSpace(chargeRef)
	result["n8n_execution_id"] = n8nExecutionID
	tier, err := s.activities.ProbeTierAtomicActivity(ctx, pb.ProbeTierActivityInput{JobId: jobID, AccountId: accountID, ProxyUrl: s.protocolProxyURL(ctx, jobID)})
	mergeActionData(result, "probe_tier", structMap(tier.GetData()))
	if err != nil {
		return nil, s.markActionFailed(ctx, jobID, stepProbeTier, jobstatus.FailedRecoverable, true, false, err, result)
	}
	if err := s.activities.MarkJobSucceededActivity(ctx, pb.JobSuccessInput{JobId: jobID, Result: structData(result)}); err != nil {
		return nil, err
	}
	s.deleteGoPayRuntimeSecrets(ctx, jobID)
	return &n8nGoPayPaymentStepResult{JobID: jobID, AccountID: accountID, N8NExecutionID: n8nExecutionID, Action: actionGoPayPayment, Step: "finish", Success: true, ChargeRef: strings.TrimSpace(chargeRef), Data: result}, nil
}

func (s *Server) FailN8NGoPayPayment(ctx context.Context, jobID string, n8nExecutionID string, flowID string, errorMessage string, data map[string]any) (any, error) {
	_ = flowID
	return s.FailN8NGoPay(ctx, actionGoPayPayment, jobID, n8nExecutionID, errorMessage, data)
}

func (s *Server) prepareN8NGoPayPayment(ctx context.Context, step string, jobID string, accountID string, n8nExecutionID string, flowID string, checkoutURL string, checkoutSessionID string, stateJSON string) (any, error) {
	jobID, accountID, n8nExecutionID = normalizeN8NGoPayQRISIDs(jobID, accountID, n8nExecutionID)
	if err := s.bindN8NGoPayExecution(ctx, jobID, n8nExecutionID); err != nil {
		return nil, err
	}
	params, err := s.jobStore.Params(ctx, jobID)
	if err != nil {
		return nil, err
	}
	gopayAccountID := goPayAppAccountID(params["gopay_account_id"])
	input := pb.GoPayActivityInput{JobId: jobID, AccountId: accountID, UseAccountToken: false, Tokenization: firstNonEmpty(params["tokenization"], "true"), GopayAccountId: gopayAccountID, StateJson: firstNonEmpty(stateJSON, "{}"), ProxyUrl: s.protocolProxyURL(ctx, jobID)}
	var out pb.GoPayPaymentPrepareOutput
	if step == stepGoPayPaymentPrepareCheckout {
		out, err = s.activities.GoPayPaymentPrepareCheckoutActivity(ctx, input)
	} else {
		input.PreparedFlowId = strings.TrimSpace(flowID)
		input.CheckoutUrl = strings.TrimSpace(checkoutURL)
		input.CheckoutSessionId = strings.TrimSpace(checkoutSessionID)
		out, err = s.activities.GoPayPaymentPrepareLinkActivity(ctx, input)
	}
	data := structMap(out.GetData())
	result := goPayPaymentPrepareResult(jobID, accountID, n8nExecutionID, step, gopayAccountID, out, data, err == nil)
	if err != nil {
		return result, s.markActionFailed(ctx, jobID, step, jobstatus.FailedRetryable, false, true, err, data)
	}
	if out.GetRetryableFreshCheckout() {
		message := strings.TrimSpace(stringMapValue(data, "error_message"))
		if message == "" {
			message = "chatgpt approve blocked"
		}
		err := fmt.Errorf("payment prepare link blocked: %s", message)
		return result, s.markActionFailed(ctx, jobID, step, jobstatus.FailedRetryable, false, true, err, data)
	}
	return result, nil
}

func goPayPaymentBaseData(params map[string]string) map[string]any {
	data := map[string]any{
		"action":              actionGoPayPayment,
		"uses_gopay_app_flow": true,
		"uses_account_token":  true,
	}
	if params != nil {
		data["source_job_id"] = params["source_job_id"]
		data["gopay_account_id"] = goPayAppAccountID(params["gopay_account_id"])
		data["tokenization"] = firstNonEmpty(params["tokenization"], "true")
	}
	return data
}

func goPayPaymentPrepareResult(jobID string, accountID string, n8nExecutionID string, step string, gopayAccountID string, out pb.GoPayPaymentPrepareOutput, data map[string]any, success bool) *n8nGoPayPaymentStepResult {
	return &n8nGoPayPaymentStepResult{JobID: jobID, AccountID: accountID, N8NExecutionID: n8nExecutionID, Action: actionGoPayPayment, Step: step, Success: success, GopayAccountID: gopayAccountID, FlowID: out.GetFlowId(), CheckoutURL: out.GetCheckoutUrl(), CheckoutSessionID: out.GetCheckoutSessionId(), UseAccountToken: out.GetUseAccountToken(), StateJSON: out.GetStateJson(), SnapToken: strings.TrimSpace(out.GetSnapToken()), SnapTokenPresent: strings.TrimSpace(out.GetSnapToken()) != "", Data: data}
}

func goPayPaymentChannelOTPCheckResult(jobID string, accountID string, n8nExecutionID string, activationID string, issuedAfterUnix int64, source string, found bool, data map[string]any) *n8nGoPayPaymentStepResult {
	if data == nil {
		data = map[string]any{}
	}
	data["otp_found"] = found
	data["otp_issued_after_unix"] = issuedAfterUnix
	if strings.TrimSpace(source) != "" {
		data["otp_source"] = strings.TrimSpace(source)
	}
	return &n8nGoPayPaymentStepResult{JobID: jobID, AccountID: accountID, N8NExecutionID: n8nExecutionID, Action: actionGoPayPayment, Step: stepGoPayPayment, Success: true, ActivationID: strings.TrimSpace(activationID), OTPFound: found, OTPSource: strings.TrimSpace(source), OTPIssuedAfterUnix: issuedAfterUnix, Data: data}
}
