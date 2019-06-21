package sobjects

type Account struct {
	BaseSObject
	BillingCity       string `force:",omitempty"`
	BillingCountry    string `force:",omitempty"`
	BillingPostalCode string `force:",omitempty"`
	BillingState      string `force:",omitempty"`
	BillingStreet     string `force:",omitempty"`
	PlatformId        string `force:"League_Platform_ID__c,omitempty"`
}

func (a Account) ApiName() string {
	return "Account"
}
