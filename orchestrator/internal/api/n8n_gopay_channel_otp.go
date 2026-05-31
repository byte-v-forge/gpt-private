//go:build private_plugins

package api

import (
	"context"
	"strings"

	"orchestrator/internal/channelotpwait"
	"orchestrator/internal/contracts"
)

const (
	goPayChannelOTPResumeSecretPrefix   = "n8n-gopay-channel-otp-resume-url:"
	goPayChannelOTPResumeSecretKeyParam = "gopay_channel_otp_resume_url_secret_key"
	channelOTPParam                     = contracts.JobParamChannelOTP
	channelOTPSubmittedAtParam          = contracts.JobParamChannelOTPSubmittedAtUnix
)

func (s *Server) AwaitN8NGoPayChannelOTP(ctx context.Context, action string, req rawN8NActionRequest) (any, error) {
	channel, target := goPayChannelOTPWaitTarget(req)
	return s.awaitN8NChannelOTP(ctx, goPayChannelOTPWaitRequest(action, req, channel, target), goPayChannelOTPWaitConfig(channel))
}

func goPayChannelOTPWaitTarget(req rawN8NActionRequest) (string, string) {
	channel := normalizeGoPayOTPChannel(firstNonEmpty(req.Channel, req.OTPChannel))
	switch channel {
	case channelotpwait.ChannelSMS:
		return channelotpwait.ChannelSMS, strings.TrimSpace(firstNonEmpty(req.Target, req.ActivationID))
	case channelotpwait.ChannelWA:
		return channelotpwait.ChannelWA, strings.TrimSpace(firstNonEmpty(req.Target, req.GopayAccountID))
	}
	if strings.TrimSpace(req.ActivationID) != "" {
		return channelotpwait.ChannelSMS, strings.TrimSpace(firstNonEmpty(req.Target, req.ActivationID))
	}
	return channelotpwait.ChannelWA, strings.TrimSpace(firstNonEmpty(req.Target, req.GopayAccountID))
}

func goPayChannelOTPWaitRequest(action string, req rawN8NActionRequest, channel string, target string) n8nChannelOTPWaitRequest {
	return n8nChannelOTPWaitRequest{
		Action:           action,
		JobID:            req.JobID,
		AccountID:        req.AccountID,
		N8NExecutionID:   req.N8NExecutionID,
		Operation:        req.Operation,
		Channel:          channel,
		Target:           target,
		StepName:         req.DataString("step_name"),
		TimeoutSeconds:   firstNonZeroInt32(req.OTPTimeoutSeconds, 300),
		OTPIssuedAfter:   req.OTPIssuedAfter,
		ResumeURL:        req.ResumeURL,
		OTPParam:         channelOTPParam,
		SubmittedAtParam: channelOTPSubmittedAtParam,
	}
}

func goPayChannelOTPWaitConfig(channel string) n8nChannelOTPWaitConfig {
	return n8nChannelOTPWaitConfig{
		Channel:              channel,
		ResumeSecretPrefix:   goPayChannelOTPResumeSecretPrefix,
		ResumeSecretKeyParam: goPayChannelOTPResumeSecretKeyParam,
	}
}

func goPayOTPChannelTarget(otpChannel string, activationID string, gopayAccountID string) (string, string) {
	if normalizeGoPayOTPChannel(otpChannel) == channelotpwait.ChannelSMS {
		return channelotpwait.ChannelSMS, strings.TrimSpace(activationID)
	}
	return channelotpwait.ChannelWA, goPayAppAccountID(gopayAccountID)
}

func goPayOTPCheckData(otpChannel string, target string, issuedAfterUnix int64) map[string]any {
	return map[string]any{
		"otp_channel":           channelotpwait.NormalizeChannel(otpChannel),
		"channel_otp_target":    channelotpwait.NormalizeTarget(otpChannel, target),
		"otp_found":             false,
		"otp_issued_after_unix": issuedAfterUnix,
		"channel_otp_converged": true,
	}
}
