package paymentsvc

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/byte-v-forge/common-lib/stringx"
)

func (c *charger) probePlusTrialCheckout(ctx context.Context) (trialProbe, error) {
	csID, err := c.createCheckout(ctx)
	if err != nil {
		return trialProbe{}, fmt.Errorf("checkout create failed: %w", err)
	}
	checkoutURL := stringx.FirstNonEmpty(c.checkoutURL, "https://checkout.stripe.com/c/pay/"+csID)
	initData, err := c.stripeInit(ctx, csID)
	if err != nil {
		return trialProbe{
			CheckoutSessionID: csID,
			CheckoutURL:       checkoutURL,
			Checked:           false,
			PlusTrialEligible: false,
			Source:            "stripe_init_error",
			ErrorMessage:      "stripe init failed: " + truncateString(err.Error(), 500),
		}, nil
	}
	amount, source := selectCheckoutAmount(initData)
	checked := source != ""
	return trialProbe{
		CheckoutSessionID: csID,
		CheckoutURL:       checkoutURL,
		Checked:           checked,
		PlusTrialEligible: checked && amount == 0,
		Amount:            amount,
		Currency:          strings.ToUpper(stringAt(initData, "currency")),
		Source:            source,
		ErrorMessage:      mapBoolString(checked, "", "stripe init did not expose checkout amount"),
	}, nil
}

func (c *charger) resolveExpectedAmount(initData map[string]any) (string, string, error) {
	if override := strings.TrimSpace(c.cfg.Runtime["expected_amount"]); override != "" {
		amount, ok := parseStripeAmount(override)
		if !ok {
			return "", "", fmt.Errorf("invalid runtime expected amount: %q", override)
		}
		return strconv.FormatInt(amount, 10), "runtime.expected_amount", nil
	}
	amount, source := selectCheckoutAmount(initData)
	if source == "" {
		if configBool(c.cfg.Runtime["fail_on_unknown_expected_amount"]) {
			return "", "", fmt.Errorf("stripe init did not expose checkout amount; refusing confirm")
		}
		return "0", "fallback_zero_unknown", nil
	}
	if amount != 0 && !configBool(c.cfg.Runtime["allow_nonzero_expected_amount"]) {
		currency := stringx.FirstNonEmpty(strings.ToUpper(stringAt(initData, "currency")), "UNKNOWN")
		return "", "", fmt.Errorf("checkout amount is %d %s from %s, not free-trial 0; refusing to confirm payment", amount, currency, source)
	}
	return strconv.FormatInt(amount, 10), source, nil
}
