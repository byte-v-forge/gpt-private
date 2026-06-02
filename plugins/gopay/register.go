package gopay

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
	payment := n8nAction(gptplugin.N8NActionSpec{
		ActionID:                actionGoPayPayment,
		DisplayName:             "GoPay Payment",
		WorkflowKey:             "gopay-payment",
		WorkflowIDPrefix:        "gopay-payment-",
		StartPath:               "/workflows/gopay-payment",
		ActionScope:             "gopay-payment",
		WebhookPath:             "gpt/gopay-payment",
		ActionPathPrefix:        "/actions/gopay-payment",
		RequestProto:            "gpt_private.GoPayPaymentRequest",
		ResponseProto:           "gpt_private.GoPayPaymentResponse",
		Button:                  actionButton("GoPay 支付", "account_detail"),
		ExtraButtons:            []gptplugin.ActionButtonSpec{actionButton("GoPay 支付", "account_row")},
		RequiredAccountStatuses: []string{gptplugin.AccountStatusRegistered},
		Capabilities:            []string{gptplugin.CapabilityPayment, gptplugin.CapabilityN8NWorkflow},
		StaleSteps:              staleSteps(),
	})
	qris := n8nAction(gptplugin.N8NActionSpec{
		ActionID:                actionGoPayQRISPaymentActivate,
		DisplayName:             "GoPay QRIS Payment Activate",
		WorkflowKey:             "gopay-qris-payment-activate",
		WorkflowIDPrefix:        "gopay-qris-payment-activate-",
		StartPath:               "/workflows/gopay-qris-payment-activate",
		ActionScope:             "gopay-qris-payment-activate",
		WebhookPath:             "gpt/gopay-qris-payment-activate",
		ActionPathPrefix:        "/actions/gopay-qris-payment-activate",
		RequestProto:            "gpt_private.GoPayQRISPaymentActivateRequest",
		ResponseProto:           "gpt_private.GoPayPaymentResponse",
		Button:                  actionButton("QRIS 激活", "account_detail"),
		RequiredAccountStatuses: []string{gptplugin.AccountStatusRegistered},
		Capabilities:            []string{gptplugin.CapabilityPayment, gptplugin.CapabilityActivation, gptplugin.CapabilityN8NWorkflow},
		StaleSteps:              staleSteps(),
	})
	return []gptplugin.ActionDefinition{
		payment,
		qris,
	}
}

func n8nAction(spec gptplugin.N8NActionSpec) gptplugin.ActionDefinition {
	spec.Owner = "gpt-private"
	spec.ActionAPIKind = gptplugin.ActionAPIKindRawN8N
	return gptplugin.BuildN8NAction(spec)
}

func actionButton(label string, placement string) gptplugin.ActionButtonSpec {
	return gptplugin.ActionButtonSpec{Label: label, Placement: placement}
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
