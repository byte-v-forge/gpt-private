package app

import (
	"context"
	"net/http"

	"github.com/byte-v-forge/common-lib/httpjson"
)

type PinWebClient struct{ domain *domainClient }

func (c *PinWebClient) TokenizePIN(ctx context.Context, pin string, challengeID string, clientID string, expected ...int) (*httpjson.Response, error) {
	return c.domain.post(ctx, "/api/v1/users/pin/tokens/nb", map[string]any{
		"challenge_id": challengeID,
		"client_id":    clientID,
		"pin":          pin,
	}, expected...)
}

func pinWebHeaderPolicy(_ *session, _ string, rawURL string, _ string, _ []byte, extra http.Header) (http.Header, error) {
	headers := pinWebHeaders()
	mergeHeaderValues(headers, extra)
	setRequestHost(headers, rawURL)
	return headers, nil
}
