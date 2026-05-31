package paymentsvc

import "github.com/byte-v-forge/common-lib/stringx"

func accessTokenTier(token string, sourcePrefix string) tierProbe {
	auth := accessTokenAuthClaims(token)
	for _, key := range []string{"chatgpt_plan_type", "chatgpt_planType", "plan_type", "planType"} {
		if result := tierResult(auth[key], sourcePrefix+"."+key); result.Checked {
			return result
		}
	}
	return tierProbe{}
}

func accessTokenAuthClaims(token string) map[string]any {
	payload := decodeJWTPayload(token)
	if payload == nil {
		return nil
	}
	auth, _ := payload["https://api.openai.com/auth"].(map[string]any)
	return auth
}

func accessTokenAccountID(token string) string {
	auth := accessTokenAuthClaims(token)
	return stringx.FirstNonEmpty(stringAt(auth, "chatgpt_account_id"), stringAt(auth, "account_id"))
}
