package paymentsvc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/byte-v-forge/common-lib/stringx"
	"github.com/google/uuid"
)

func (c *charger) stripeConfirm(ctx context.Context, csID, pmID string) (map[string]any, error) {
	initData, err := c.stripeInit(ctx, csID)
	if err != nil {
		return nil, err
	}
	expectedAmount, _, err := c.resolveExpectedAmount(initData)
	if err != nil {
		return nil, err
	}
	chatGPTReturn := "https://chatgpt.com/checkout/verify?stripe_session_id=" + csID + "&processor_entity=" + c.processorEntityOrDefault() + "&plan_type=plus"
	returnURL := "https://checkout.stripe.com/c/pay/" + csID + "?returned_from_redirect=true&ui_mode=custom&return_url=" + url.QueryEscape(chatGPTReturn)
	clientSessionID := uuid.NewString()
	form := url.Values{
		"guid":                                   {uuidHex()},
		"muid":                                   {uuidHex()},
		"sid":                                    {uuidHex()},
		"payment_method":                         {pmID},
		"init_checksum":                          {stringAt(initData, "init_checksum")},
		"version":                                {stringx.FirstNonEmpty(c.cfg.Runtime["version"], "fed52f3bc6")},
		"expected_amount":                        {expectedAmount},
		"expected_payment_method_type":           {"gopay"},
		"return_url":                             {returnURL},
		"elements_session_client[session_id]":    {newElementsSessionID()},
		"elements_session_client[locale]":        {"en"},
		"elements_session_client[referrer_host]": {"chatgpt.com"},
		"elements_session_client[is_aggregation_expected]":                         {"false"},
		"elements_session_client[elements_init_source]":                            {"custom_checkout"},
		"elements_session_client[client_betas][0]":                                 {"custom_checkout_server_updates_1"},
		"elements_session_client[client_betas][1]":                                 {"custom_checkout_manual_approval_1"},
		"client_attribution_metadata[client_session_id]":                           {clientSessionID},
		"client_attribution_metadata[checkout_session_id]":                         {csID},
		"client_attribution_metadata[merchant_integration_source]":                 {"checkout"},
		"client_attribution_metadata[merchant_integration_subtype]":                {"payment-element"},
		"client_attribution_metadata[merchant_integration_version]":                {"custom"},
		"client_attribution_metadata[payment_intent_creation_flow]":                {"deferred"},
		"client_attribution_metadata[payment_method_selection_flow]":               {"automatic"},
		"client_attribution_metadata[merchant_integration_additional_elements][0]": {"payment"},
		"client_attribution_metadata[merchant_integration_additional_elements][1]": {"address"},
		"consent[terms_of_service]":                                                {"accepted"},
		"key":                                                                      {c.cfg.StripePublishableKey},
		"_stripe_version":                                                          {stripeCheckoutVersion},
	}
	if value := strings.TrimSpace(c.cfg.Runtime["js_checksum"]); value != "" {
		form.Set("js_checksum", value)
	}
	if value := strings.TrimSpace(c.cfg.Runtime["rv_timestamp"]); value != "" {
		form.Set("rv_timestamp", value)
	}
	resp, err := c.stripe.request(ctx, http.MethodPost, "https://api.stripe.com/v1/payment_pages/"+csID+"/confirm", requestOptions{formBody: form})
	if err != nil {
		return nil, err
	}
	if resp.status == http.StatusBadRequest && strings.Contains(strings.ToLower(string(resp.body)), "terms of service") {
		form.Set("consent[terms_of_service]", "accepted")
		resp, err = c.stripe.request(ctx, http.MethodPost, "https://api.stripe.com/v1/payment_pages/"+csID+"/confirm", requestOptions{formBody: form})
		if err != nil {
			return nil, err
		}
	}
	if resp.status != http.StatusOK {
		return nil, fmt.Errorf("stripe confirm %d: %s", resp.status, resp.excerpt(500))
	}
	return resp.json, nil
}
