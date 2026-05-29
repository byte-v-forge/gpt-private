package privateflows

import (
	"github.com/byte-v-forge/gpt-private/plugins/internal/plugincatalog"
	"github.com/byte-v-forge/gpt/pkg/gptplugin"
)

func configSchema() gptplugin.ConfigSchema {
	return gptplugin.ConfigSchema{
		PluginKey:   "private_flows",
		DisplayName: "Private GPT Flows",
		Owner:       "gpt-private",
		Fields: []gptplugin.ConfigField{
			configField("registration_otp_timeout_seconds", "Registration OTP Timeout Seconds", gptplugin.ConfigFieldDurationSeconds, "180"),
		},
	}
}

func configField(key string, label string, kind gptplugin.ConfigFieldKind, defaultValue string) gptplugin.ConfigField {
	return plugincatalog.Field(key, label, kind, defaultValue)
}
