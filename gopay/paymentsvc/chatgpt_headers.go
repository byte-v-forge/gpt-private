package paymentsvc

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const gptAcceptLanguage = "en-US,en;q=0.9"

func (c *GptClient) applyChatGPTHeaders() {
	if c == nil || c.session == nil {
		return
	}
	mergeHeader(c.session.headers, c.chatGPTHeaders(nil))
}

func (c *GptClient) chatGPTHeaders(extra http.Header) http.Header {
	fingerprint := c.fingerprint()
	headers := http.Header{}
	applyChatGPTBrowserHeaders(headers, fingerprint)
	mergeHeader(headers, extra)
	if headers.Get("Accept") == "" {
		headers.Set("Accept", "*/*")
	}
	if headers.Get("Accept-Language") == "" {
		headers.Set("Accept-Language", gptAcceptLanguage)
	}
	headers.Set("Origin", "https://chatgpt.com")
	if headers.Get("Referer") == "" {
		headers.Set("Referer", "https://chatgpt.com/")
	}
	if !hasExplicitHeader(extra, "sec-ch-ua") {
		headers.Del("sec-ch-ua")
	}
	if !hasExplicitHeader(extra, "sec-ch-ua-mobile") {
		headers.Del("sec-ch-ua-mobile")
	}
	if !hasExplicitHeader(extra, "sec-ch-ua-platform") {
		headers.Del("sec-ch-ua-platform")
	}
	if !hasExplicitHeader(extra, "oai-language") {
		headers.Del("oai-language")
	}
	if !hasExplicitHeader(extra, "oai-device-id") {
		headers.Del("oai-device-id")
	}
	if !hasExplicitHeader(extra, "sec-fetch-dest") {
		headers.Del("sec-fetch-dest")
	}
	if !hasExplicitHeader(extra, "sec-fetch-mode") {
		headers.Del("sec-fetch-mode")
	}
	if !hasExplicitHeader(extra, "sec-fetch-site") {
		headers.Del("sec-fetch-site")
	}
	return headers
}

func hasExplicitHeader(headers http.Header, name string) bool {
	if len(headers) == 0 {
		return false
	}
	if _, ok := headers[http.CanonicalHeaderKey(name)]; ok {
		return true
	}
	for key := range headers {
		if strings.EqualFold(key, name) {
			return true
		}
	}
	return false
}

func applyChatGPTBrowserHeaders(headers http.Header, fingerprint browserFingerprint) {
	if headers == nil {
		return
	}
	if fingerprint.UserAgent != "" {
		headers.Set("User-Agent", fingerprint.UserAgent)
	}
	if fingerprint.AcceptLanguage != "" {
		headers.Set("Accept-Language", fingerprint.AcceptLanguage)
	}
}

func requireGptOpenAIURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	if host == "chatgpt.com" || strings.HasSuffix(host, ".chatgpt.com") || host == "openai.com" || strings.HasSuffix(host, ".openai.com") {
		return nil
	}
	return fmt.Errorf("gpt client refuses non gpt/openai url: %s", host)
}
