package privateflows

import (
	"github.com/byte-v-forge/gpt/pkg/gptplugin"
)

func configSchema() gptplugin.ConfigSchema {
	return gptplugin.ConfigSchema{
		PluginKey:   "private_flows",
		DisplayName: "Private GPT Flows",
		Owner:       "gpt-private",
		Fields: []gptplugin.ConfigField{
			gptplugin.Field("registration_otp_timeout_seconds", "Registration OTP Timeout Seconds", gptplugin.ConfigFieldDurationSeconds, "180"),
		},
	}
}
