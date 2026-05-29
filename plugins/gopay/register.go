package gopay

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
	app := withCapabilities(
		n8nAction(actionGoPayApp, "GoPay App", "gopay-app", "gopay-app-", "/workflows/gopay-app", "gopay-app", "gpt/gopay-app", "/actions/gopay-app", "orchestrator.GoPayAppRequest", "orchestrator.GoPayAppResponse", "GoPay App", "gopay"),
		capabilityGoPay,
		gptplugin.CapabilityN8NWorkflow,
	)
	payment := withCapabilities(
		plugincatalog.WithUIButton(plugincatalog.WithRequiredStatuses(n8nAction(actionGoPayPayment, "GoPay Payment", "gopay-payment", "gopay-payment-", "/workflows/gopay-payment", "gopay-payment", "gpt/gopay-payment", "/actions/gopay-payment", "orchestrator.GoPayPaymentRequest", "orchestrator.GoPayPaymentResponse", "GoPay 支付", "account_detail"), "REGISTERED"), "GoPay 支付", "account_row"),
		capabilityGoPay,
		gptplugin.CapabilityPayment,
		gptplugin.CapabilityN8NWorkflow,
	)
	qris := withCapabilities(
		plugincatalog.WithRequiredStatuses(n8nAction(actionGoPayQRISPaymentActivate, "GoPay QRIS Payment Activate", "gopay-qris-payment-activate", "gopay-qris-payment-activate-", "/workflows/gopay-qris-payment-activate", "gopay-qris-payment-activate", "gpt/gopay-qris-payment-activate", "/actions/gopay-qris-payment-activate", "orchestrator.GoPayQRISPaymentActivateRequest", "orchestrator.GoPayPaymentResponse", "QRIS 激活", "account_detail"), "REGISTERED"),
		capabilityGoPay,
		gptplugin.CapabilityPayment,
		gptplugin.CapabilityActivation,
		gptplugin.CapabilityN8NWorkflow,
	)
	wa := withCapabilities(
		plugincatalog.WithRequiredStatuses(n8nAction(actionGoPayWAPayment, "GoPay WA Payment", "gopay-wa-payment", "gopay-wa-payment-", "/workflows/gopay-wa-payment", "gopay-wa-payment", "gpt/gopay-wa-payment", "/actions/gopay-wa-payment", "orchestrator.GoPayWAPaymentRequest", "orchestrator.GoPayPaymentResponse", "WA 支付", "account_detail"), "REGISTERED"),
		capabilityGoPay,
		gptplugin.CapabilityPayment,
		gptplugin.CapabilityN8NWorkflow,
	)
	rebind := withCapabilities(
		n8nAction(actionGoPayPaymentRebind, "GoPay Payment Rebind", "gopay-payment-rebind", "gopay-payment-rebind-", "/workflows/gopay-payment-rebind", "gopay-payment-rebind", "gpt/gopay-payment-rebind", "/actions/gopay-payment-rebind", "orchestrator.GoPayPaymentRebindRequest", "orchestrator.GoPayPaymentResponse", "重绑支付", "account_detail"),
		capabilityGoPay,
		gptplugin.CapabilityPayment,
		gptplugin.CapabilityN8NWorkflow,
	)
	return []gptplugin.ActionDefinition{
		app,
		payment,
		qris,
		wa,
		rebind,
	}
}

func withCapabilities(def gptplugin.ActionDefinition, capabilities ...string) gptplugin.ActionDefinition {
	return plugincatalog.WithCapabilities(def, capabilities...)
}

func n8nAction(actionID string, displayName string, workflowKey string, workflowIDPrefix string, startPath string, actionScope string, webhookPath string, actionPathPrefix string, requestProto string, responseProto string, buttonLabel string, placement string) gptplugin.ActionDefinition {
	def := plugincatalog.N8NAction("gpt-private", actionID, displayName, workflowKey, workflowIDPrefix, startPath, actionScope, webhookPath, actionPathPrefix, requestProto, responseProto, buttonLabel, placement, staleSteps())
	def.Workflow.ActionAPIKind = gptplugin.ActionAPIKindRawN8N
	return def
}

func staleSteps() []string {
	return []string{
		"claimed",
		stepGoPayPaymentPrepare,
		stepGoPayPaymentPrepareCheckout,
		stepGoPayPaymentPrepareRefresh,
		stepGoPayPaymentPrepareLink,
		stepGoPayPayment,
		stepGoPayResolveWAPhone,
		stepGoPayAppEnsureBalance,
		stepGoPayAppEnsureBalanceConfirm,
		stepGoPayAppLogin,
		stepGoPayAppEnsurePINSetup,
		stepGoPayAppChangePhone,
		stepGoPayAppChangePhoneGetNumber,
		stepGoPayAppChangePhoneStart,
		stepGoPayAppChangePhoneSMSWait,
		stepGoPayAppChangePhoneRetry,
		stepGoPayAppChangePhoneCancel,
		stepGoPayAppChangePhoneComplete,
		stepGoPayAppSMSFinish,
		stepGoPayAppSMSRequestMore,
		stepGoPayAppSignup,
		stepGoPayAppSignupRetry,
		stepGoPayAppSignupPhoneCancel,
		stepGoPayAppSignupPhone,
		stepGoPayAppGenerateDeviceProxy,
		stepGoPayAppCheckPhone,
		stepGoPayAppDeactivate,
		stepGoPayAppDeactivateStart,
		stepGoPayAppDeactivateSMSWait,
		stepGoPayAppDeactivateComplete,
		stepGoPayAppStatus,
	}
}
