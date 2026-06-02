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
	payment := withCapabilities(
		plugincatalog.WithUIButton(plugincatalog.WithRequiredStatuses(n8nAction(actionGoPayPayment, "GoPay Payment", "gopay-payment", "gopay-payment-", "/workflows/gopay-payment", "gopay-payment", "gpt/gopay-payment", "/actions/gopay-payment", "gpt_private.GoPayPaymentRequest", "gpt_private.GoPayPaymentResponse", "GoPay 支付", "account_detail"), "REGISTERED"), "GoPay 支付", "account_row"),
		gptplugin.CapabilityPayment,
		gptplugin.CapabilityN8NWorkflow,
	)
	qris := withCapabilities(
		plugincatalog.WithRequiredStatuses(n8nAction(actionGoPayQRISPaymentActivate, "GoPay QRIS Payment Activate", "gopay-qris-payment-activate", "gopay-qris-payment-activate-", "/workflows/gopay-qris-payment-activate", "gopay-qris-payment-activate", "gpt/gopay-qris-payment-activate", "/actions/gopay-qris-payment-activate", "gpt_private.GoPayQRISPaymentActivateRequest", "gpt_private.GoPayPaymentResponse", "QRIS 激活", "account_detail"), "REGISTERED"),
		gptplugin.CapabilityPayment,
		gptplugin.CapabilityActivation,
		gptplugin.CapabilityN8NWorkflow,
	)
	return []gptplugin.ActionDefinition{
		payment,
		qris,
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
