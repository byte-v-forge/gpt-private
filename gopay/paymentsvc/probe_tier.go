package paymentsvc

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	paidTierRE          = regexp.MustCompile(`(^|[_:/\-\s])(chatgpt[_:/\-\s]*)?(plus|pro|team|business|enterprise)([_:/\-\s]|$)`)
	freeTierRE          = regexp.MustCompile(`(^|[_:/\-\s])(free|none|anonymous|unauthenticated)([_:/\-\s]|$)`)
	paidMarkerCleanerRE = regexp.MustCompile(`[^a-z0-9]+`)
)

func tierResult(value any, source string) tierProbe {
	tier := normalizeTier(value)
	if tier == "" {
		return tierProbe{}
	}
	return tierProbe{Checked: true, PlusActive: tier != "free", PlanType: tier, Tier: tier, Source: source}
}

func normalizeTier(value any) string {
	text := strings.ToLower(strings.TrimSpace(fmt.Sprint(value)))
	if text == "" || text == "<nil>" {
		return ""
	}
	for _, tier := range []string{"plus", "pro", "team", "business", "enterprise", "free"} {
		if text == tier {
			return tier
		}
	}
	if match := paidTierRE.FindStringSubmatch(text); len(match) > 3 {
		return match[3]
	}
	if freeTierRE.MatchString(text) {
		return "free"
	}
	return ""
}

func detectPlusActiveFromSessionPayload(payload map[string]any) tierProbe {
	if payload == nil {
		return tierProbe{Source: "auth_session"}
	}
	if result := accessTokenTier(stringAt(payload, "accessToken"), "accessToken.auth"); result.Checked {
		return result
	}
	if account, ok := payload["account"].(map[string]any); ok {
		for _, key := range []string{"planType", "plan_type", "accountPlan", "account_plan", "tier", "plan"} {
			if result := tierResult(account[key], "account."+key); result.Checked {
				return result
			}
		}
	}
	if walkPaidMarker(payload) {
		return tierProbe{Checked: true, PlusActive: true, PlanType: "paid", Tier: "paid", Source: "auth_session:paid_marker"}
	}
	checked := payload["user"] != nil || payload["accessToken"] != nil
	return tierProbe{Checked: checked, PlusActive: false, PlanType: "free", Tier: "free", Source: "auth_session:no_paid_marker"}
}

func walkPaidMarker(value any) bool {
	switch typed := value.(type) {
	case map[string]any:
		for key, item := range typed {
			clean := strings.ToLower(paidMarkerCleanerRE.ReplaceAllString(key, ""))
			if b, ok := item.(bool); ok && b {
				for _, marker := range []string{"haspaid", "haspaidsubscription", "hasactivepaidsubscription", "isplus", "ispaid", "subscribed"} {
					if clean == marker {
						return true
					}
				}
			}
			if result := tierResult(item, key); result.Checked && result.PlusActive {
				return true
			}
			if walkPaidMarker(item) {
				return true
			}
		}
	case []any:
		for _, item := range typed {
			if walkPaidMarker(item) {
				return true
			}
		}
	}
	return false
}
