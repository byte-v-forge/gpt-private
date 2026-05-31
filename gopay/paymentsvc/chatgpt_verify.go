package paymentsvc

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

func (c *charger) chatGPTVerify(ctx context.Context, csID string) map[string]any {
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := c.cs.request(ctx, http.MethodGet, "https://chatgpt.com/checkout/verify", requestOptions{
			query: url.Values{"stripe_session_id": []string{csID}, "processor_entity": []string{c.processorEntityOrDefault()}, "plan_type": []string{"plus"}},
		})
		if err == nil && resp.status == http.StatusOK {
			return map[string]any{"state": "succeeded", "cs_id": csID}
		}
		select {
		case <-ctx.Done():
			return map[string]any{"state": "verify_timeout", "cs_id": csID}
		case <-time.After(2 * time.Second):
		}
	}
	return map[string]any{"state": "verify_timeout", "cs_id": csID}
}
