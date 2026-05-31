//go:build private_plugins

package contracts

const (
	ActionGoPayApp                 = "GOPAY_APP"
	ActionGoPayPayment             = "GOPAY_PAYMENT"
	ActionGoPayQRISPaymentActivate = "GOPAY_QRIS_PAYMENT_ACTIVATE"
	ActionGoPayWAPayment           = "GOPAY_WA_PAYMENT"
	ActionGoPayPaymentRebind       = "GOPAY_PAYMENT_REBIND"
)

const (
	CapabilityGoPay = "gopay"
)

const (
	ManualGoPayPaymentConfirmationParam = "manual_gopay_payment_confirmed"
)

const (
	StepGoPayAppLogin                = "gopay_app_ensure_token_available"
	StepGoPayAppChangePhone          = "gopay_app_change_phone"
	StepGoPayAppChangePhoneGetNumber = "gopay_app_change_phone_get_number"
	StepGoPayAppChangePhoneStart     = "gopay_app_change_phone_start"
	StepGoPayAppChangePhoneSMSWait   = "gopay_app_change_phone_sms_wait"
	StepGoPayAppChangePhoneRetry     = "gopay_app_change_phone_retry"
	StepGoPayAppChangePhoneCancel    = "gopay_app_change_phone_cancel"
	StepGoPayAppChangePhoneComplete  = "gopay_app_change_phone_complete"
	StepGoPayAppSignupPhone          = "gopay_app_signup_phone"
	StepGoPayAppGenerateDeviceProxy  = "gopay_app_generate_device_proxy"
	StepGoPayAppCheckPhone           = "gopay_app_check_phone"
	StepGoPayResolveWAPhone          = "gopay_resolve_wa_phone"
	StepGoPayAppDeactivate           = "gopay_app_deactivate"
	StepGoPayAppDeactivateStart      = "gopay_app_deactivate_start"
	StepGoPayAppDeactivateSMSWait    = "gopay_app_deactivate_sms_wait"
	StepGoPayAppDeactivateSMSFinish  = "gopay_app_deactivate_sms_finish"
	StepGoPayAppDeactivateComplete   = "gopay_app_deactivate_complete"
	StepGoPayAppSignup               = "gopay_app_signup"
	StepGoPayAppSignupRetry          = "gopay_app_signup_retry"
	StepGoPayAppSignupPhoneCancel    = "gopay_app_signup_phone_cancel"
	StepGoPayAppStatus               = "gopay_app_status"
	StepGoPayAppEnsurePINSetup       = "gopay_app_ensure_pin_setup"
	StepGoPayAppSMSFinish            = "gopay_app_sms_finish"
	StepGoPayAppSMSRequestMore       = "gopay_app_sms_request_more"
	StepGoPayPaymentPrepare          = "gopay_payment_prepare"
	StepGoPayPaymentPrepareCheckout  = "gopay_payment_prepare_checkout"
	StepGoPayPaymentPrepareRefresh   = "gopay_payment_prepare_checkout_refresh"
	StepGoPayPaymentPrepareLink      = "gopay_payment_prepare_link"
	StepGoPayPayment                 = "gopay_payment"
)
