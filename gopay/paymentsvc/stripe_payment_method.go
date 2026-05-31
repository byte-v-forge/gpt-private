package paymentsvc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/byte-v-forge/common-lib/browserfingerprint"
	"github.com/byte-v-forge/common-lib/stringx"
	"github.com/google/uuid"
)

const stripeCheckoutVersion = "2025-03-31.basil; checkout_server_update_beta=v1; checkout_manual_approval_preview=v1"

func (c *charger) stripeCreatePaymentMethod(ctx context.Context, csID string) (string, error) {
	billing := c.cfg.Billing
	runtimeVersion := stringx.FirstNonEmpty(c.cfg.Runtime["version"], "fed52f3bc6")
	clientSessionID := uuid.NewString()
	form := url.Values{
		"billing_details[name]":                                                    {stringx.FirstNonEmpty(billing["name"], "John Doe")},
		"billing_details[email]":                                                   {stringx.FirstNonEmpty(billing["email"], "buyer@example.com")},
		"billing_details[address][country]":                                        {stringx.FirstNonEmpty(billing["country"], "US")},
		"billing_details[address][line1]":                                          {stringx.FirstNonEmpty(billing["line1"], "3110 Sunset Boulevard")},
		"billing_details[address][city]":                                           {stringx.FirstNonEmpty(billing["city"], "Los Angeles")},
		"billing_details[address][postal_code]":                                    {stringx.FirstNonEmpty(billing["postal_code"], "90026")},
		"billing_details[address][state]":                                          {stringx.FirstNonEmpty(billing["state"], "CA")},
		"type":                                                                     {"gopay"},
		"payment_user_agent":                                                       {fmt.Sprintf("stripe.js/%s; stripe-js-v3/%s; payment-element; deferred-intent", runtimeVersion, runtimeVersion)},
		"referrer":                                                                 {"https://chatgpt.com"},
		"time_on_page":                                                             {fmt.Sprintf("%d", 25000+browserfingerprint.RandomIndex(30001))},
		"client_attribution_metadata[client_session_id]":                           {clientSessionID},
		"client_attribution_metadata[checkout_session_id]":                         {csID},
		"client_attribution_metadata[elements_session_id]":                         {newElementsSessionID()},
		"client_attribution_metadata[elements_session_config_id]":                  {uuid.NewString()},
		"client_attribution_metadata[merchant_integration_source]":                 {"elements"},
		"client_attribution_metadata[merchant_integration_subtype]":                {"payment-element"},
		"client_attribution_metadata[merchant_integration_version]":                {"2021"},
		"client_attribution_metadata[payment_intent_creation_flow]":                {"deferred"},
		"client_attribution_metadata[payment_method_selection_flow]":               {"automatic"},
		"client_attribution_metadata[merchant_integration_additional_elements][0]": {"payment"},
		"client_attribution_metadata[merchant_integration_additional_elements][1]": {"address"},
		"guid":            {uuidHex()},
		"muid":            {uuidHex()},
		"sid":             {uuidHex()},
		"key":             {c.cfg.StripePublishableKey},
		"_stripe_version": {stripeCheckoutVersion},
	}
	resp, err := c.stripe.request(ctx, http.MethodPost, "https://api.stripe.com/v1/payment_methods", requestOptions{formBody: form})
	if err != nil {
		return "", err
	}
	if err := resp.require(http.StatusOK, "stripe payment_methods"); err != nil {
		return "", err
	}
	pmID := stringAt(resp.json, "id")
	if !strings.HasPrefix(pmID, "pm_") {
		return "", fmt.Errorf("stripe payment_methods: bad response %s", resp.excerpt(300))
	}
	return pmID, nil
}
