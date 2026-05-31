package paymentsvc

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/byte-v-forge/common-lib/stringx"
)

func (c *charger) createCheckout(ctx context.Context) (string, error) {
	if err := c.cs.setProxy(c.checkoutProfile.ProxyURL); err != nil {
		return "", fmt.Errorf("checkout proxy switch: %w", err)
	}
	c.chatGPTWarmup(ctx)
	plan := c.cfg.CheckoutPlan
	checkoutMode := stringx.FirstNonEmpty(plan["checkout_ui_mode"], "custom")
	body := map[string]any{
		"entry_point": stringx.FirstNonEmpty(plan["entry_point"], "all_plans_pricing_modal"),
		"plan_name":   stringx.FirstNonEmpty(plan["plan_name"], "chatgptplusplan"),
		"billing_details": map[string]any{
			"country":  stringx.FirstNonEmpty(plan["billing_country"], "ID"),
			"currency": stringx.FirstNonEmpty(plan["billing_currency"], "IDR"),
		},
		"checkout_ui_mode": checkoutMode,
	}
	body["cancel_url"] = stringx.FirstNonEmpty(plan["cancel_url"], "https://chatgpt.com/#pricing")
	if promo := stringx.FirstNonEmpty(plan["promo_campaign_id"], "plus-1-month-free"); promo != "" {
		body["promo_campaign"] = map[string]any{"promo_campaign_id": promo, "is_coupon_from_query_param": false}
	}
	resp, err := c.cs.request(ctx, http.MethodPost, "https://chatgpt.com/backend-api/payments/checkout", requestOptions{
		jsonBody: body,
		headers:  c.chatGPTCheckoutHeaders(),
	})
	if err != nil {
		return "", err
	}
	if resp.status >= 400 {
		return "", fmt.Errorf("checkout create failed: status=%d %s", resp.status, resp.excerpt(500))
	}
	csID := extractCheckoutSessionID(resp.json)
	if csID == "" {
		return "", fmt.Errorf("checkout create: bad response %s", jsonExcerpt(resp.json, 500))
	}
	c.checkoutURL = checkoutURLFromResponse(resp.json, csID)
	c.processorEntity = stringx.FirstNonEmpty(extractProcessorEntity(resp.json), c.processorEntity, "openai_llc")
	return csID, nil
}

func (c *charger) chatGPTWarmup(ctx context.Context) {
	billingCountry := strings.ToUpper(stringx.FirstNonEmpty(c.cfg.CheckoutPlan["billing_country"], "ID"))
	warmups := []struct {
		method string
		url    string
		accept string
		body   any
	}{
		{http.MethodGet, "https://chatgpt.com/", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8", nil},
		{http.MethodGet, "https://chatgpt.com/api/auth/session", "application/json", nil},
		{http.MethodGet, "https://chatgpt.com/backend-api/accounts/check/v4-2023-04-27?timezone_offset_min=-540", "application/json", nil},
		{http.MethodGet, "https://chatgpt.com/backend-api/accounts/domain-density-eligibility", "application/json", nil},
		{http.MethodGet, "https://chatgpt.com/backend-api/checkout_pricing_config/countries", "application/json", nil},
		{http.MethodGet, "https://chatgpt.com/backend-api/checkout_pricing_config/configs/" + billingCountry, "application/json", nil},
	}
	for _, item := range warmups {
		select {
		case <-ctx.Done():
			return
		default:
		}
		resp, _ := c.cs.request(ctx, item.method, item.url, requestOptions{
			headers:  c.chatGPTWarmupHeaders(item.accept, item.url),
			jsonBody: item.body,
		})
		if resp != nil && item.url == "https://chatgpt.com/api/auth/session" && resp.status == http.StatusOK {
			if accessToken := stringAt(resp.json, "accessToken"); accessToken != "" {
				c.cs.setAccessToken(accessToken)
			}
		}
	}
	if cookie := mergeCookieHeaders(c.cs.header("Cookie"), c.cs.cookieHeader("https://chatgpt.com/")); cookie != "" {
		c.cs.setHeader("Cookie", cookie)
	}
}

func (c *charger) processorEntityOrDefault() string {
	return stringx.FirstNonEmpty(c.processorEntity, extractProcessorEntityFromURL(c.checkoutURL), "openai_llc")
}

func (c *charger) checkoutApprovalURL(csID string) string {
	return "https://chatgpt.com/checkout/" + c.processorEntityOrDefault() + "/" + csID
}
