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
			configField("change_phone_max_failures", "Change Phone Max Failures", gptplugin.ConfigFieldInteger, "3"),
			configField("change_phone_disabled", "Change Phone Disabled", gptplugin.ConfigFieldBoolean, "false"),
			configField("change_phone_otp_retry_attempts", "Change Phone OTP Retry Attempts", gptplugin.ConfigFieldInteger, "1"),
			configField("change_phone_get_number_retry_seconds", "Change Phone Get Number Retry Seconds", gptplugin.ConfigFieldDurationSeconds, "5"),
			configField("sms_route_min_available_count", "SMS Min Available Count", gptplugin.ConfigFieldInteger, "10"),
			configField("sms_route_limit", "SMS Route Limit", gptplugin.ConfigFieldInteger, "10"),
			configField("sms_route_max_price_amount", "SMS Max Price Amount", gptplugin.ConfigFieldString, ""),
			configFieldHelp("sms_route_failure_scope_key", "SMS Failure Scope Key", gptplugin.ConfigFieldString, "gopay", "按业务隔离 SMS 路由临时禁用名单；同一 scope 内统计同一路由失败。"),
			configFieldHelp("sms_route_failure_threshold", "SMS Failure Threshold", gptplugin.ConfigFieldInteger, "3", "失败统计窗口内连续失败达到该次数后，该 SMS 路由进入临时禁用名单。"),
			configFieldHelp("sms_route_failure_window_seconds", "SMS Failure Window Seconds", gptplugin.ConfigFieldDurationSeconds, "600", "失败次数统计窗口，单位秒；例如 600 表示最近 10 分钟。"),
			configFieldHelp("sms_route_disable_ttl_seconds", "SMS Disable TTL Seconds", gptplugin.ConfigFieldDurationSeconds, "600", "路由进入临时禁用名单后的 TTL，单位秒。"),
		},
	}
}

func configField(key string, label string, kind gptplugin.ConfigFieldKind, defaultValue string) gptplugin.ConfigField {
	return plugincatalog.Field(key, label, kind, defaultValue)
}

func configFieldHelp(key string, label string, kind gptplugin.ConfigFieldKind, defaultValue string, helpText string) gptplugin.ConfigField {
	field := configField(key, label, kind, defaultValue)
	field.HelpText = helpText
	return field
}
