//go:build private_plugins

package api

import "strings"

func normalizeGoPayOTPChannel(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "wa", "whatsapp", "otp_wa":
		return "wa"
	case "sms", "otp_sms":
		return "sms"
	default:
		return ""
	}
}
