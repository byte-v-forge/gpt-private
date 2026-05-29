package gopay

import (
	"github.com/byte-v-forge/gpt-private/plugins/internal/plugincatalog"
	"github.com/byte-v-forge/gpt/pkg/gptplugin"
)

func configSchema() gptplugin.ConfigSchema {
	return gptplugin.ConfigSchema{
		PluginKey:   "gopay",
		DisplayName: "GoPay",
		Owner:       "gpt-private",
		Fields: []gptplugin.ConfigField{
			configField("otp_timeout_seconds", "OTP Timeout Seconds", gptplugin.ConfigFieldDurationSeconds, "180"),
			configField("app_step_body_limit", "App Step Body Limit", gptplugin.ConfigFieldInteger, "6000"),
			configField("app_link_payment_timeout_seconds", "Link Payment Timeout Seconds", gptplugin.ConfigFieldDurationSeconds, "180"),
			configField("app_unlink_timeout_seconds", "Unlink Timeout Seconds", gptplugin.ConfigFieldDurationSeconds, "15"),
			configField("add_balance_envelope_link", "Envelope Link", gptplugin.ConfigFieldSecret, ""),
			configField("add_balance_transfer_instructions", "Transfer Instructions", gptplugin.ConfigFieldSecret, ""),
			configField("add_balance_transfer_amount_rp", "Transfer Amount Rp", gptplugin.ConfigFieldInteger, "1"),
			configField("add_balance_transfer_currency", "Transfer Currency", gptplugin.ConfigFieldString, "IDR"),
			configField("add_balance_confirm_timeout_seconds", "Add Balance Confirm Timeout Seconds", gptplugin.ConfigFieldDurationSeconds, "1800"),
			configField("change_phone_max_failures", "Change Phone Max Failures", gptplugin.ConfigFieldInteger, "3"),
			configField("change_phone_disabled", "Change Phone Disabled", gptplugin.ConfigFieldBoolean, "false"),
			configField("change_phone_otp_retry_attempts", "Change Phone OTP Retry Attempts", gptplugin.ConfigFieldInteger, "1"),
			configField("change_phone_get_number_retry_seconds", "Change Phone Get Number Retry Seconds", gptplugin.ConfigFieldDurationSeconds, "5"),
		},
	}
}

func configField(key string, label string, kind gptplugin.ConfigFieldKind, defaultValue string) gptplugin.ConfigField {
	return plugincatalog.Field(key, label, kind, defaultValue)
}
