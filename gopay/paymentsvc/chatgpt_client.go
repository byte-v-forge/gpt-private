package paymentsvc

import (
	"context"
	"fmt"
	"strings"
)

type GptClient struct {
	session *httpSession
	profile requestProfile
}

func (s *Server) newGptClient(ctx context.Context, cred credential, profile requestProfile) (*GptClient, error) {
	if cred.empty() {
		return nil, fmt.Errorf("auth missing: need session_token or access_token")
	}
	profile = profile.withDefaults(defaultRequestProfile(profile.Name))
	fingerprint := profile.fingerprint()
	session, err := newHTTPSession(profile.ProxyURL, fingerprint)
	if err != nil {
		return nil, err
	}
	client := &GptClient{session: session, profile: profile}
	client.applyChatGPTHeaders()
	if cred.accessToken != "" {
		client.setAccessToken(cred.accessToken)
	}
	if cookie := chatGPTCookieHeader(cred.sessionToken, ""); cookie != "" {
		client.setHeader("Cookie", cookie)
	}
	return client, nil
}

func newGptClientWithFingerprint(proxyURL string, fingerprint browserFingerprint) (*GptClient, error) {
	session, err := newHTTPSession(proxyURL, fingerprint)
	if err != nil {
		return nil, err
	}
	client := &GptClient{session: session, profile: requestProfile{ProxyURL: proxyURL, TLSProfile: fingerprint.TLSProfileName, Locale: fingerprint.OAILanguage}}
	client.applyChatGPTHeaders()
	return client, nil
}

func (c *GptClient) request(ctx context.Context, method, rawURL string, opts requestOptions) (*httpResult, error) {
	if c == nil || c.session == nil {
		return nil, fmt.Errorf("gpt client is nil")
	}
	if err := requireGptOpenAIURL(rawURL); err != nil {
		return nil, err
	}
	c.applyChatGPTHeaders()
	opts.headers = c.chatGPTHeaders(opts.headers)
	return c.session.request(ctx, method, rawURL, opts)
}

func (c *GptClient) setProxy(proxyURL string) error {
	if c == nil || c.session == nil {
		return fmt.Errorf("gpt client is nil")
	}
	return c.session.setProxy(proxyURL)
}

func (c *GptClient) close() {
	if c != nil && c.session != nil {
		c.session.close()
	}
}

func (c *GptClient) cookieHeader(rawURL string) string {
	if c == nil || c.session == nil {
		return ""
	}
	return c.session.cookieHeader(rawURL)
}

func (c *GptClient) fingerprint() browserFingerprint {
	if c != nil && c.session != nil {
		return c.session.fingerprint.withFallback(c.profile.Locale)
	}
	return defaultRequestProfile("checkout").fingerprint()
}

func (c *GptClient) header(name string) string {
	if c == nil || c.session == nil {
		return ""
	}
	return c.session.headers.Get(name)
}

func (c *GptClient) setHeader(name string, value string) {
	if c == nil || c.session == nil {
		return
	}
	c.session.headers.Set(name, value)
}

func (c *GptClient) setAuthorization(value string) {
	c.setHeader("Authorization", value)
}

func (c *GptClient) setAccessToken(token string) {
	token = strings.TrimSpace(token)
	if token == "" {
		return
	}
	c.setAuthorization("Bearer " + token)
	if accountID := accessTokenAccountID(token); accountID != "" {
		c.setHeader("ChatGPT-Account-Id", accountID)
	}
}
