package paymentsvc

type tierProbe struct {
	Checked      bool
	PlusActive   bool
	PlanType     string
	Tier         string
	Source       string
	ErrorMessage string
}

type trialProbe struct {
	Checked           bool
	PlusTrialEligible bool
	PlusActive        bool
	PlanType          string
	Amount            int64
	Currency          string
	Source            string
	CheckoutURL       string
	CheckoutSessionID string
	ErrorMessage      string
}
