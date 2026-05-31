//go:build private_plugins

package activities

import "orchestrator/internal/contracts"

const (
	actionGoPayApp                 = contracts.ActionGoPayApp
	actionGoPayPayment             = contracts.ActionGoPayPayment
	actionGoPayQRISPaymentActivate = contracts.ActionGoPayQRISPaymentActivate
	actionGoPayWAPayment           = contracts.ActionGoPayWAPayment
	actionGoPayPaymentRebind       = contracts.ActionGoPayPaymentRebind

	stepGoPayAppLogin                = contracts.StepGoPayAppLogin
	stepGoPayAppChangePhone          = contracts.StepGoPayAppChangePhone
	stepGoPayAppChangePhoneGetNumber = contracts.StepGoPayAppChangePhoneGetNumber
	stepGoPayAppChangePhoneStart     = contracts.StepGoPayAppChangePhoneStart
	stepGoPayAppChangePhoneSMSWait   = contracts.StepGoPayAppChangePhoneSMSWait
	stepGoPayAppChangePhoneRetry     = contracts.StepGoPayAppChangePhoneRetry
	stepGoPayAppChangePhoneCancel    = contracts.StepGoPayAppChangePhoneCancel
	stepGoPayAppChangePhoneComplete  = contracts.StepGoPayAppChangePhoneComplete
	stepGoPayAppSignupPhone          = contracts.StepGoPayAppSignupPhone
	stepGoPayAppGenerateDeviceProxy  = contracts.StepGoPayAppGenerateDeviceProxy
	stepGoPayAppCheckPhone           = contracts.StepGoPayAppCheckPhone
	stepGoPayResolveWAPhone          = contracts.StepGoPayResolveWAPhone
	stepGoPayAppDeactivate           = contracts.StepGoPayAppDeactivate
	stepGoPayAppDeactivateStart      = contracts.StepGoPayAppDeactivateStart
	stepGoPayAppDeactivateSMSWait    = contracts.StepGoPayAppDeactivateSMSWait
	stepGoPayAppDeactivateSMSFinish  = contracts.StepGoPayAppDeactivateSMSFinish
	stepGoPayAppDeactivateComplete   = contracts.StepGoPayAppDeactivateComplete
	stepGoPayAppSignup               = contracts.StepGoPayAppSignup
	stepGoPayAppSignupRetry          = contracts.StepGoPayAppSignupRetry
	stepGoPayAppSignupPhoneCancel    = contracts.StepGoPayAppSignupPhoneCancel
	stepGoPayAppStatus               = contracts.StepGoPayAppStatus
	stepGoPayAppEnsurePINSetup       = contracts.StepGoPayAppEnsurePINSetup
	stepGoPayAppSMSFinish            = contracts.StepGoPayAppSMSFinish
	stepGoPayAppSMSRequestMore       = contracts.StepGoPayAppSMSRequestMore
	stepGoPayPaymentPrepare          = contracts.StepGoPayPaymentPrepare
	stepGoPayPaymentPrepareCheckout  = contracts.StepGoPayPaymentPrepareCheckout
	stepGoPayPaymentPrepareRefresh   = contracts.StepGoPayPaymentPrepareRefresh
	stepGoPayPaymentPrepareLink      = contracts.StepGoPayPaymentPrepareLink
	stepGoPayPayment                 = contracts.StepGoPayPayment
)
