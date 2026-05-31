package paymentsvc

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type StripeClient struct {
	session *httpSession
}

func newStripeClient(session *httpSession) (*StripeClient, error) {
	if session == nil {
		return nil, fmt.Errorf("stripe client requires http session")
	}
	return &StripeClient{session: session}, nil
}

func (c *StripeClient) request(ctx context.Context, method, rawURL string, opts requestOptions) (*httpResult, error) {
	if c == nil || c.session == nil {
		return nil, fmt.Errorf("stripe client is nil")
	}
	if err := requireStripeURL(rawURL); err != nil {
		return nil, err
	}
	return c.session.request(ctx, method, rawURL, opts)
}

func (c *StripeClient) requestIfStripe(ctx context.Context, method, rawURL string, opts requestOptions) (*httpResult, bool, error) {
	if isStripeURL(rawURL) {
		resp, err := c.request(ctx, method, rawURL, opts)
		return resp, true, err
	}
	return nil, false, nil
}

func requireStripeURL(rawURL string) error {
	if isStripeURL(rawURL) {
		return nil
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	return fmt.Errorf("stripe client refuses non stripe url: %s", strings.ToLower(strings.TrimSpace(parsed.Hostname())))
}

func isStripeURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	return host == "stripe.com" || strings.HasSuffix(host, ".stripe.com")
}
