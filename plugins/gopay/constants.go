package gopay

const (
	actionGoPayPayment             = "GOPAY_PAYMENT"
	actionGoPayQRISPaymentActivate = "GOPAY_QRIS_PAYMENT_ACTIVATE"
	actionGoPayWAPayment           = "GOPAY_WA_PAYMENT"
	actionGoPayPaymentRebind       = "GOPAY_PAYMENT_REBIND"
)

const (
	stepGoPayAppLogin                = "gopay_app_ensure_token_available"
	stepGoPayAppChangePhone          = "gopay_app_change_phone"
	stepGoPayAppChangePhoneGetNumber = "gopay_app_change_phone_get_number"
	stepGoPayAppChangePhoneStart     = "gopay_app_change_phone_start"
	stepGoPayAppChangePhoneSMSWait   = "gopay_app_change_phone_sms_wait"
	stepGoPayAppChangePhoneRetry     = "gopay_app_change_phone_retry"
	stepGoPayAppChangePhoneCancel    = "gopay_app_change_phone_cancel"
	stepGoPayAppChangePhoneComplete  = "gopay_app_change_phone_complete"
	stepGoPayAppSignupPhone          = "gopay_app_signup_phone"
	stepGoPayAppGenerateDeviceProxy  = "gopay_app_generate_device_proxy"
	stepGoPayAppCheckPhone           = "gopay_app_check_phone"
	stepGoPayResolveWAPhone          = "gopay_resolve_wa_phone"
	stepGoPayAppDeactivate           = "gopay_app_deactivate"
	stepGoPayAppDeactivateStart      = "gopay_app_deactivate_start"
	stepGoPayAppDeactivateSMSWait    = "gopay_app_deactivate_sms_wait"
	stepGoPayAppDeactivateComplete   = "gopay_app_deactivate_complete"
	stepGoPayAppSignup               = "gopay_app_signup"
	stepGoPayAppSignupRetry          = "gopay_app_signup_retry"
	stepGoPayAppSignupPhoneCancel    = "gopay_app_signup_phone_cancel"
	stepGoPayAppStatus               = "gopay_app_status"
	stepGoPayAppEnsurePINSetup       = "gopay_app_ensure_pin_setup"
	stepGoPayAppSMSFinish            = "gopay_app_sms_finish"
	stepGoPayAppSMSRequestMore       = "gopay_app_sms_request_more"
	stepGoPayPaymentPrepare          = "gopay_payment_prepare"
	stepGoPayPaymentPrepareCheckout  = "gopay_payment_prepare_checkout"
	stepGoPayPaymentPrepareRefresh   = "gopay_payment_prepare_checkout_refresh"
	stepGoPayPaymentPrepareLink      = "gopay_payment_prepare_link"
	stepGoPayPayment                 = "gopay_payment"
)
