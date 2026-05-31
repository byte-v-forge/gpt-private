package paymentsvc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/google/uuid"
)

var midtransSnapRedirectRE = regexp.MustCompile(`app\.midtrans\.com/snap/v[14]/redirection/([a-f0-9-]{36})`)

func (c *charger) fetchPMRedirectSnapToken(ctx context.Context, pmURL string) (string, error) {
	if token := regexpSnapToken(pmURL); token != "" {
		return token, nil
	}
	resp, usedStripe, err := c.stripe.requestIfStripe(ctx, http.MethodGet, pmURL, requestOptions{noRedirect: true})
	if err != nil {
		return "", err
	}
	if !usedStripe {
		resp, err = c.paymentHTTP.request(ctx, http.MethodGet, pmURL, requestOptions{noRedirect: true})
		if err != nil {
			return "", err
		}
	}
	if resp.status < 300 || resp.status > 399 {
		return "", fmt.Errorf("pm-redirects: expected redirect, got %d", resp.status)
	}
	location := resp.headers.Get("Location")
	if token := regexpSnapToken(location); token != "" {
		return token, nil
	}
	return "", fmt.Errorf("pm-redirects: no midtrans token in redirect")
}

func (c *charger) followRedirectToMidtrans(ctx context.Context, csID string) (string, error) {
	deadline := time.Now().Add(60 * time.Second)
	lastErr := ""
	for time.Now().Before(deadline) {
		query := url.Values{
			"elements_session_client[client_betas][0]":                        {"custom_checkout_server_updates_1"},
			"elements_session_client[client_betas][1]":                        {"custom_checkout_manual_approval_1"},
			"elements_session_client[elements_init_source]":                   {"custom_checkout"},
			"elements_session_client[referrer_host]":                          {"chatgpt.com"},
			"elements_session_client[session_id]":                             {newElementsSessionID()},
			"elements_session_client[stripe_js_id]":                           {uuid.NewString()},
			"elements_session_client[locale]":                                 {"en"},
			"elements_session_client[is_aggregation_expected]":                {"false"},
			"elements_options_client[saved_payment_method][enable_save]":      {"never"},
			"elements_options_client[saved_payment_method][enable_redisplay]": {"never"},
			"key":             {c.cfg.StripePublishableKey},
			"_stripe_version": {stripeCheckoutVersion},
		}
		resp, err := c.stripe.request(ctx, http.MethodGet, "https://api.stripe.com/v1/payment_pages/"+csID, requestOptions{query: query})
		if err != nil {
			lastErr = err.Error()
		} else if resp.status == http.StatusOK {
			if pmURL := extractRedirectToURL(resp.json); pmURL != "" {
				return c.fetchPMRedirectSnapToken(ctx, pmURL)
			}
			lastErr = fmt.Sprintf("setup_intent status=%q payment_intent status=%q invoice.payment_intent status=%q payment_status=%q status=%q",
				stringAt(resp.json, "setup_intent", "status"),
				stringAt(resp.json, "payment_intent", "status"),
				stringAt(resp.json, "invoice", "payment_intent", "status"),
				stringAt(resp.json, "payment_status"),
				stringAt(resp.json, "status"),
			)
		} else {
			lastErr = fmt.Sprintf("http %d: %s", resp.status, resp.excerpt(150))
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(time.Second):
		}
	}
	return "", fmt.Errorf("snap_token resolution timeout: %s", lastErr)
}

func regexpSnapToken(value string) string {
	match := midtransSnapRedirectRE.FindStringSubmatch(value)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
