package paymentsvc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/byte-v-forge/common-lib/stringx"
)

func (c *charger) prepareUntilLinking(ctx context.Context, checkoutSessionID, checkoutURL string) (map[string]any, error) {
	csID, state, err := c.prepareCheckout(ctx, checkoutSessionID, checkoutURL, 1)
	if err != nil {
		return nil, err
	}
	if stringAt(state, "checkout_supplied") == "true" {
		return c.prepareCheckoutSessionUntilLinking(ctx, csID)
	}
	var lastErr error
	for attempt := int(intAt(state, "checkout_attempt")); attempt <= 2; attempt++ {
		if attempt > int(intAt(state, "checkout_attempt")) {
			var refreshErr error
			csID, state, refreshErr = c.prepareCheckout(ctx, "", "", attempt)
			if refreshErr != nil {
				return nil, refreshErr
			}
		}
		prepared, err := c.prepareCheckoutSessionUntilLinking(ctx, csID)
		if err == nil {
			prepared["checkout_attempt"] = attempt
			return prepared, nil
		}
		lastErr = err
		if !isChatGPTApproveBlocked(err) {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(2+attempt) * time.Second):
		}
	}
	return nil, lastErr
}

func (c *charger) prepareCheckout(ctx context.Context, checkoutSessionID, checkoutURL string, attempt int) (string, map[string]any, error) {
	checkoutURL = strings.TrimSpace(checkoutURL)
	csID := strings.TrimSpace(checkoutSessionID)
	if csID == "" && checkoutURL != "" {
		csID = extractCheckoutSessionID(map[string]any{"url": checkoutURL})
	}
	if checkoutURL != "" && csID == "" {
		return "", nil, fmt.Errorf("checkout_url does not contain checkout_session_id")
	}
	if csID != "" {
		if !strings.HasPrefix(csID, "cs_") {
			return "", nil, fmt.Errorf("invalid checkout_session_id: %s", csID)
		}
		c.checkoutURL = stringx.FirstNonEmpty(checkoutURL, "https://checkout.stripe.com/c/pay/"+csID)
		c.processorEntity = stringx.FirstNonEmpty(extractProcessorEntityFromURL(checkoutURL), c.processorEntity, "openai_llc")
		return csID, map[string]any{
			"state":             "checkout",
			"cs_id":             csID,
			"processor_entity":  c.processorEntityOrDefault(),
			"checkout_url":      c.checkoutURL,
			"stripe_pk":         c.cfg.StripePublishableKey,
			"checkout_attempt":  attempt,
			"checkout_supplied": "true",
		}, nil
	}

	csID, err := c.createCheckout(ctx)
	if err != nil {
		return "", nil, err
	}
	return csID, map[string]any{
		"state":             "checkout",
		"cs_id":             csID,
		"processor_entity":  c.processorEntityOrDefault(),
		"checkout_url":      c.checkoutURL,
		"stripe_pk":         c.cfg.StripePublishableKey,
		"checkout_attempt":  attempt,
		"checkout_supplied": "false",
	}, nil
}

func (c *charger) prepareCheckoutSessionUntilLinking(ctx context.Context, csID string) (map[string]any, error) {
	pmID, err := c.stripeCreatePaymentMethod(ctx, csID)
	if err != nil {
		return nil, err
	}
	confirmData, err := c.stripeConfirm(ctx, csID, pmID)
	if err != nil {
		return nil, err
	}
	redirectURL := extractRedirectToURL(confirmData)
	var snapToken string
	if redirectURL != "" {
		snapToken, err = c.fetchPMRedirectSnapToken(ctx, redirectURL)
	} else {
		if err = c.chatGPTApprove(ctx, csID); err != nil {
			return nil, err
		}
		snapToken, err = c.followRedirectToMidtrans(ctx, csID)
	}
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"state":            "prepared",
		"cs_id":            csID,
		"processor_entity": c.processorEntityOrDefault(),
		"checkout_url":     c.checkoutURL,
		"stripe_pk":        c.cfg.StripePublishableKey,
		"snap_token":       snapToken,
	}, nil
}
