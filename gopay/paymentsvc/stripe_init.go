package paymentsvc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

func (c *charger) stripeInit(ctx context.Context, csID string) (map[string]any, error) {
	form := url.Values{
		"browser_locale":                                   {"en-US"},
		"browser_timezone":                                 {"Asia/Shanghai"},
		"elements_session_client[client_betas][0]":         {"custom_checkout_server_updates_1"},
		"elements_session_client[client_betas][1]":         {"custom_checkout_manual_approval_1"},
		"elements_session_client[elements_init_source]":    {"custom_checkout"},
		"elements_session_client[referrer_host]":           {"chatgpt.com"},
		"elements_session_client[stripe_js_id]":            {uuid.NewString()},
		"elements_session_client[locale]":                  {"en"},
		"elements_session_client[is_aggregation_expected]": {"false"},
		"key":             {c.cfg.StripePublishableKey},
		"_stripe_version": {stripeCheckoutVersion},
	}
	resp, err := c.stripe.request(ctx, http.MethodPost, "https://api.stripe.com/v1/payment_pages/"+csID+"/init", requestOptions{formBody: form})
	if err != nil {
		return nil, err
	}
	if err := resp.require(http.StatusOK, "stripe init"); err != nil {
		return nil, err
	}
	if !containsString(resp.json["payment_method_types"], "gopay") {
		return nil, fmt.Errorf("checkout does not support GoPay: currency=%s payment_method_types=%v", stringAt(resp.json, "currency"), resp.json["payment_method_types"])
	}
	if stringAt(resp.json, "init_checksum") == "" {
		return nil, fmt.Errorf("stripe init: no init_checksum %s", resp.excerpt(300))
	}
	return resp.json, nil
}
