package sobjects

type Account struct {
	BaseSObject
	BillingCity       string `force:",omitempty"`
	BillingCountry    string `force:",omitempty"`
	BillingPostalCode string `force:",omitempty"`
	BillingState      string `force:",omitempty"`
	BillingStreet     string `force:",omitempty"`
}

func (a Account) ApiName() string {
	return "Account"
}
