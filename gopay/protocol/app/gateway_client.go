package app

import (
	"context"
	"net/http"

	"github.com/byte-v-forge/common-lib/httpjson"
)

type GatewayClient struct{ domain *domainClient }

func (c *GatewayClient) Get(ctx context.Context, path string, expected ...int) (*httpjson.Response, error) {
	return c.domain.get(ctx, path, expected...)
}

func (c *GatewayClient) Post(ctx context.Context, path string, body any, expected ...int) (*httpjson.Response, error) {
	return c.domain.post(ctx, path, body, expected...)
}

func gatewayHeaderPolicy(_ *session, _ string, rawURL string, _ string, body []byte, extra http.Header) (http.Header, error) {
	headers := gopayGatewayHeaders(len(body) > 0)
	mergeHeaderValues(headers, extra)
	setRequestHost(headers, rawURL)
	return headers, nil
}
