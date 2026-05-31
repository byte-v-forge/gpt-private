package paymentsvc

import (
	"context"
	"strings"
)

type charger struct {
	cfg              Config
	checkoutProfile  requestProfile
	paymentProfile   requestProfile
	cs               *GptClient
	paymentHTTP      *httpSession
	stripe           *StripeClient
	countryCode      string
	phone            string
	pin              string
	tokenization     string
	checkoutURL      string
	processorEntity  string
	midtransMerchant string
}

func (s *Server) newCharger(ctx context.Context, cred credential, phone, countryCode, pin, tokenization string) (*charger, error) {
	checkoutProfile := cred.chatGPTProfile(s.cfg.CheckoutProfile)
	paymentProfile := s.cfg.PaymentProfile
	paymentFingerprint := paymentProfile.fingerprint()
	cs, err := s.newGptClient(ctx, cred, checkoutProfile)
	if err != nil {
		return nil, err
	}
	paymentHTTP, err := newHTTPSession(paymentProfile.ProxyURL, paymentFingerprint)
	if err != nil {
		cs.close()
		return nil, err
	}
	paymentFingerprint.applyBrowserHeaders(paymentHTTP.headers)
	stripe, err := newStripeClient(paymentHTTP)
	if err != nil {
		cs.close()
		paymentHTTP.close()
		return nil, err
	}
	return &charger{
		cfg:             s.cfg,
		checkoutProfile: checkoutProfile,
		paymentProfile:  paymentProfile,
		cs:              cs,
		paymentHTTP:     paymentHTTP,
		stripe:          stripe,
		countryCode:     normalizeCountryCode(countryCode),
		phone:           normalizeDigits(phone),
		pin:             strings.TrimSpace(pin),
		tokenization:    normalizeTokenization(tokenization),
	}, nil
}

func (c *charger) close() {
	if c == nil {
		return
	}
	if c.cs != nil {
		c.cs.close()
	}
	if c.paymentHTTP != nil {
		c.paymentHTTP.close()
	}
}

func (c *charger) requiresManualConfirmation() bool {
	return requiresManualPaymentConfirmation(c.tokenization)
}

func requiresManualPaymentConfirmation(tokenization string) bool {
	value := strings.ToLower(strings.TrimSpace(tokenization))
	return value == "false" || value == "qris"
}

func isQRISTokenization(tokenization string) bool {
	return strings.EqualFold(strings.TrimSpace(tokenization), "qris")
}

func normalizeTokenization(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultTokenization
	}
	return value
}
