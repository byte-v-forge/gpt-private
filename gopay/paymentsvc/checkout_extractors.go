package paymentsvc

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	checkoutSessionIDRE       = regexp.MustCompile(`cs_(?:live|test)_[A-Za-z0-9]+`)
	processorEntityCheckoutRE = regexp.MustCompile(`/checkout/([^/?#]+)/cs_(?:live|test)_[A-Za-z0-9]+`)
)

func extractCheckoutSessionID(data map[string]any) string {
	for _, key := range []string{"checkout_session_id", "session_id", "id"} {
		value := stringAt(data, key)
		if strings.HasPrefix(value, "cs_") {
			return value
		}
	}
	for _, key := range []string{"url", "stripe_hosted_url", "checkout_url"} {
		if match := checkoutSessionIDRE.FindString(stringAt(data, key)); match != "" {
			return match
		}
	}
	return ""
}

func extractProcessorEntity(data map[string]any) string {
	for _, key := range []string{"processor_entity", "processorEntity", "merchant_entity"} {
		if value := strings.TrimSpace(stringAt(data, key)); value != "" {
			return value
		}
	}
	for _, key := range []string{"url", "stripe_hosted_url", "checkout_url", "success_url", "cancel_url", "return_url"} {
		if value := extractProcessorEntityFromURL(stringAt(data, key)); value != "" {
			return value
		}
	}
	return ""
}

func extractProcessorEntityFromURL(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if parsed, err := url.Parse(value); err == nil {
		if entity := strings.TrimSpace(parsed.Query().Get("processor_entity")); entity != "" {
			return entity
		}
	}
	matches := processorEntityCheckoutRE.FindStringSubmatch(value)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func checkoutURLFromResponse(data map[string]any, csID string) string {
	for _, key := range []string{"url", "stripe_hosted_url", "checkout_url"} {
		if value := stringAt(data, key); value != "" {
			return value
		}
	}
	if csID != "" {
		return "https://checkout.stripe.com/c/pay/" + csID
	}
	return ""
}
