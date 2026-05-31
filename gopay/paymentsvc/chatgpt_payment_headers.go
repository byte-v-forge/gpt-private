package paymentsvc

import (
	"net/http"
	"strings"

	"github.com/byte-v-forge/common-lib/stringx"
	"github.com/google/uuid"
)

func (c *charger) chatGPTAuthHeaders(referer string) http.Header {
	headers := c.cs.chatGPTHeaders(http.Header{
		"Accept":       []string{"*/*"},
		"Content-Type": []string{"application/json"},
		"Origin":       []string{"https://chatgpt.com"},
		"Referer":      []string{stringx.FirstNonEmpty(referer, "https://chatgpt.com/")},
	})
	c.addChatGPTAuthMaterial(headers)
	return headers
}

func (c *charger) chatGPTWarmupHeaders(accept, rawURL string) http.Header {
	headers := http.Header{
		"Accept":  []string{accept},
		"Referer": []string{"https://chatgpt.com/"},
	}
	if strings.Contains(rawURL, "/backend-api/") {
		headers.Set("Origin", "https://chatgpt.com")
	}
	c.addChatGPTAuthMaterial(headers)
	if strings.Contains(rawURL, "/backend-api/") {
		c.addChatGPTBrowserStateHeaders(headers)
	}
	return headers
}

func (c *charger) chatGPTCheckoutHeaders() http.Header {
	headers := c.chatGPTAuthHeaders("https://chatgpt.com/")
	c.addChatGPTBrowserStateHeaders(headers)
	headers.Set("oai-session-id", uuid.NewString())
	headers.Set("x-openai-target-path", "/backend-api/payments/checkout")
	headers.Set("x-openai-target-route", "/backend-api/payments/checkout")
	return headers
}

func (c *charger) addChatGPTAuthMaterial(headers http.Header) {
	if headers == nil {
		return
	}
	if value := c.cs.header("Authorization"); value != "" {
		headers.Set("Authorization", value)
	}
	if cookie := mergeCookieHeaders(c.cs.header("Cookie"), c.cs.cookieHeader("https://chatgpt.com/")); cookie != "" {
		headers.Set("Cookie", cookie)
	}
}

func (c *charger) addChatGPTBrowserStateHeaders(headers http.Header) {
	if headers == nil {
		return
	}
	fp := c.chatGPTFingerprint()
	if fp.DeviceID != "" {
		headers.Set("oai-device-id", fp.DeviceID)
		if cookie := headers.Get("Cookie"); cookie != "" {
			headers.Set("Cookie", cookieHeaderWithDeviceID(splitCookieHeader(cookie), fp.DeviceID))
		}
	}
	headers.Set("oai-language", stringx.FirstNonEmpty(fp.OAILanguage, fp.Language, "en-US"))
	if fp.SecCHUA != "" {
		headers.Set("sec-ch-ua", fp.SecCHUA)
		headers.Set("sec-ch-ua-mobile", "?0")
	}
	if fp.SecCHPlatform != "" {
		headers.Set("sec-ch-ua-platform", fp.SecCHPlatform)
	}
	headers.Set("sec-fetch-dest", "empty")
	headers.Set("sec-fetch-mode", "cors")
	headers.Set("sec-fetch-site", "same-origin")
}

func (c *charger) chatGPTApproveHeaders(csID string) http.Header {
	// nb-register style: keep the same ChatGPT browser identity that created the
	// checkout, then move the session to the payment proxy only for approve.
	// Switching to a fresh payment-profile device can make approve return blocked.
	headers := c.chatGPTAuthHeaders(c.checkoutApprovalURL(csID))
	c.addChatGPTBrowserStateHeaders(headers)
	headers.Set("x-openai-target-path", "/backend-api/payments/checkout/approve")
	headers.Set("x-openai-target-route", "/backend-api/payments/checkout/approve")
	return headers
}

func (c *charger) chatGPTFingerprint() browserFingerprint {
	if c != nil && c.cs != nil {
		return c.cs.fingerprint()
	}
	return defaultRequestProfile("checkout").fingerprint()
}
