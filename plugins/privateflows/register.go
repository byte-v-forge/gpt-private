package privateflows

import (
	"github.com/byte-v-forge/gpt/pkg/gptplugin"
)

func Plugin() gptplugin.Plugin {
	return gptplugin.PluginFunc(Register)
}

func Register(registry gptplugin.ActionRegistry) error {
	if err := registry.RegisterActions(actions()...); err != nil {
		return err
	}
	return registry.RegisterPluginConfigs(configSchema())
}

func actions() []gptplugin.ActionDefinition {
	register := n8nAction(gptplugin.N8NActionSpec{
		ActionID:               actionRegister,
		DisplayName:            "Register",
		WorkflowKey:            "register",
		WorkflowIDPrefix:       "register-",
		StartPath:              "/workflows/register",
		ActionScope:            "register",
		WebhookPath:            "gpt/register",
		ActionPathPrefix:       "/actions/register",
		RequestProto:           "orchestrator.RegisterAccountRequest",
		ResponseProto:          "orchestrator.RegisterAccountResponse",
		Button:                 actionButton("浏览器注册", "account_header_browser"),
		RequiredFields:         []string{"email", "password"},
		BlockedAccountStatuses: []string{gptplugin.AccountStatusRegistered, gptplugin.AccountStatusActivated},
		Capabilities:           []string{gptplugin.CapabilityRegistration, gptplugin.CapabilityBrowserAuth, gptplugin.CapabilityN8NWorkflow},
		StaleSteps:             staleSteps(),
	})
	registerProtocol := n8nAction(gptplugin.N8NActionSpec{
		ActionID:               actionRegisterProtocol,
		DisplayName:            "Register Protocol",
		WorkflowKey:            "register-protocol",
		WorkflowIDPrefix:       "register-protocol-",
		StartPath:              "/workflows/register-protocol",
		ActionScope:            "register-protocol",
		WebhookPath:            "gpt/register-protocol",
		ActionPathPrefix:       "/actions/register-protocol",
		RequestProto:           "orchestrator.RegisterAccountRequest",
		ResponseProto:          "orchestrator.RegisterAccountResponse",
		Button:                 actionButton("协议注册", "account_header_protocol"),
		ExtraButtons:           []gptplugin.ActionButtonSpec{actionButton("注册", "account_row")},
		RequiredFields:         []string{"email", "password"},
		BlockedAccountStatuses: []string{gptplugin.AccountStatusRegistered, gptplugin.AccountStatusActivated},
		Capabilities:           []string{gptplugin.CapabilityRegistration, gptplugin.CapabilityProtocolAuth, gptplugin.CapabilityN8NWorkflow},
		StaleSteps:             staleSteps(),
	})
	return []gptplugin.ActionDefinition{
		register,
		registerProtocol,
	}
}

func n8nAction(spec gptplugin.N8NActionSpec) gptplugin.ActionDefinition {
	spec.Owner = "gpt-private"
	return gptplugin.BuildN8NAction(spec)
}

func actionButton(label string, placement string) gptplugin.ActionButtonSpec {
	return gptplugin.ActionButtonSpec{Label: label, Placement: placement}
}

func staleSteps() []string {
	return []string{
		"claimed",
		stepRegisterAccountStart,
		stepRegisterAccountBrowser,
		stepRegisterAccountOTPWait,
		stepRegisterAccountComplete,
		gptplugin.StepProtocolUseProxy,
		stepRegisterAccountProtocolStart,
		stepRegisterAccountProtocol,
		stepRegisterAccountProtocolOTPWait,
		stepRegisterAccountProtocolComplete,
	}
}
