package privateflows

import (
	"github.com/byte-v-forge/gpt-private/plugins/internal/plugincatalog"
	"github.com/byte-v-forge/gpt/pkg/gptplugin"
)

func init() {
	gptplugin.Register(gptplugin.PluginFunc(Register))
}

func Register(registry gptplugin.ActionRegistry) error {
	if err := registry.RegisterActions(actions()...); err != nil {
		return err
	}
	return registry.RegisterPluginConfigs(configSchema())
}

func actions() []gptplugin.ActionDefinition {
	register := withCapabilities(
		plugincatalog.WithBlockedStatuses(plugincatalog.WithRequiredFields(n8nAction(actionRegister, "Register", "register", "register-", "/workflows/register", "register", "gpt/register", "/actions/register", "orchestrator.RegisterAccountRequest", "orchestrator.RegisterAccountResponse", "浏览器注册", "account_header_browser"), "email", "password"), "REGISTERED", "ACTIVATED"),
		gptplugin.CapabilityRegistration,
		gptplugin.CapabilityBrowserAuth,
		gptplugin.CapabilityN8NWorkflow,
	)
	registerProtocol := withCapabilities(
		plugincatalog.WithUIButton(plugincatalog.WithBlockedStatuses(plugincatalog.WithRequiredFields(n8nAction(actionRegisterProtocol, "Register Protocol", "register-protocol", "register-protocol-", "/workflows/register-protocol", "register-protocol", "gpt/register-protocol", "/actions/register-protocol", "orchestrator.RegisterAccountRequest", "orchestrator.RegisterAccountResponse", "协议注册", "account_header_protocol"), "email", "password"), "REGISTERED", "ACTIVATED"), "注册", "account_row"),
		gptplugin.CapabilityRegistration,
		gptplugin.CapabilityProtocolAuth,
		gptplugin.CapabilityN8NWorkflow,
	)
	return []gptplugin.ActionDefinition{
		register,
		registerProtocol,
	}
}

func withCapabilities(def gptplugin.ActionDefinition, capabilities ...string) gptplugin.ActionDefinition {
	return plugincatalog.WithCapabilities(def, capabilities...)
}

func n8nAction(actionID string, displayName string, workflowKey string, workflowIDPrefix string, startPath string, actionScope string, webhookPath string, actionPathPrefix string, requestProto string, responseProto string, buttonLabel string, placement string) gptplugin.ActionDefinition {
	return plugincatalog.N8NAction("gpt-private", actionID, displayName, workflowKey, workflowIDPrefix, startPath, actionScope, webhookPath, actionPathPrefix, requestProto, responseProto, buttonLabel, placement, staleSteps())
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
