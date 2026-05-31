//go:build private_plugins

package api

import (
	"context"
	"strings"
)

const goPayDynamicProxyCountryCode = "ID"

func (s *Server) N8NGoPayDynamicProxySettings(ctx context.Context, jobID string, accountID string, n8nExecutionID string, purpose string) (any, error) {
	jobID, accountID, n8nExecutionID = normalizeN8NGoPayQRISIDs(jobID, accountID, n8nExecutionID)
	if err := s.bindN8NGoPayExecution(ctx, jobID, n8nExecutionID); err != nil {
		return nil, err
	}
	return s.n8nDynamicProxySettingsForGeo(ctx, jobID, accountID, n8nExecutionID, goPayDynamicProxyProfile(purpose), goPayDynamicProxyCountryCode, "", "gopay_dynamic_proxy")
}

func (s *Server) RecordN8NGoPayDynamicProxy(ctx context.Context, jobID string, accountID string, n8nExecutionID string, proxyURL string, data map[string]any, purpose string) (any, error) {
	_ = data
	return s.recordN8NDynamicProxy(ctx, jobID, accountID, n8nExecutionID, proxyURL, nil, goPayDynamicProxyProfile(purpose))
}

func (s *Server) FailN8NGoPayDynamicProxy(ctx context.Context, jobID string, accountID string, n8nExecutionID string, errorMessage string, data map[string]any) (any, error) {
	_ = data
	return s.failN8NDynamicProxy(ctx, jobID, accountID, n8nExecutionID, errorMessage, nil, goPayDynamicProxyProfile(actionGoPayPayment))
}

func goPayDynamicProxyProfile(purpose string) n8nDynamicProxyProfile {
	return n8nDynamicProxyProfile{Purpose: strings.TrimSpace(purpose)}
}
