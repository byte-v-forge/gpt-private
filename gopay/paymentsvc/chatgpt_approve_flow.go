package paymentsvc

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (c *charger) chatGPTApprove(ctx context.Context, csID string) error {
	if err := c.cs.setProxy(c.paymentProfile.ProxyURL); err != nil {
		return fmt.Errorf("chatgpt approve proxy switch: %w", err)
	}
	headers := c.chatGPTApproveHeaders(csID)
	c.chatGPTSentinelPing(ctx, c.cs)

	var lastStatus int
	var lastBody string
	for attempt := 1; attempt <= 3; attempt++ {
		resp, err := c.cs.request(ctx, http.MethodPost, "https://chatgpt.com/backend-api/payments/checkout/approve", requestOptions{
			jsonBody: map[string]any{"checkout_session_id": csID, "processor_entity": c.processorEntityOrDefault()},
			headers:  headers,
		})
		if err != nil {
			return err
		}
		lastStatus = resp.status
		lastBody = resp.excerpt(500)
		result := stringAt(resp.json, "result")
		if resp.status == http.StatusOK && result == "approved" {
			return nil
		}
		if result == "blocked" && attempt < 3 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(2+attempt) * time.Second):
			}
			c.chatGPTSentinelPing(ctx, c.cs)
			continue
		}
		if resp.status == http.StatusForbidden && strings.Contains(strings.ToLower(lastBody), "<html") {
			return fmt.Errorf("chatgpt approve cloudflare challenge 403: %s", lastBody)
		}
		if resp.status != http.StatusOK {
			return fmt.Errorf("chatgpt approve %d: %s", resp.status, lastBody)
		}
		if result == "blocked" {
			return chatGPTApproveBlockedError{status: lastStatus, body: lastBody}
		}
		return fmt.Errorf("chatgpt approve: result=%q body=%s", result, lastBody)
	}
	return chatGPTApproveBlockedError{status: lastStatus, body: lastBody}
}

func (c *charger) chatGPTSentinelPing(ctx context.Context, session *GptClient) {
	if session == nil {
		return
	}
	_, _ = session.request(ctx, http.MethodPost, "https://chatgpt.com/backend-api/sentinel/ping", requestOptions{
		jsonBody: map[string]any{},
		headers:  c.chatGPTAuthHeaders("https://chatgpt.com/"),
	})
}
