package force

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestComposite_RequestAndResponse(t *testing.T) {
	forceAPI := createTest()

	newCompositeRequests := &CompositeRequests{}

	contactID := "0032O000002KnrEQAS" //Replace this with a contact ID that exists :D

	newCompositeRequests.AllOrNone = true

	newApplication := &Application{
		ClearmatchCustomer: contactID,
		ApplicationNumber:  "A1111113",
		LoanAmount:         1111.11,
	}
	newRequest := &CompositeRequest{
		Method:      http.MethodPost,
		URL:         "/services/data/v42.0/sobjects/Applications__c",
		Body:        *newApplication,
		ReferenceID: "test",
	}

	newCompositeRequests.Add(newRequest)

	newApplication2 := &Application{
		ClearmatchCustomer: contactID,
		ApplicationNumber:  "A1111114",
		LoanAmount:         2222.22,
	}
	newRequest2 := &CompositeRequest{
		Method:      http.MethodPost,
		URL:         "/services/data/v42.0/sobjects/Applications__c",
		Body:        *newApplication2,
		ReferenceID: "test2",
	}

	newCompositeRequests.Add(newRequest2)

	resp, err := forceAPI.PostCompositeRequests(newCompositeRequests)

	st, _ := json.Marshal(resp)
	fmt.Println(string(st))

	if err != nil {
		t.Error("Invalid Composite Requests")
	}

	for _, r := range resp.CompositeResponse {
		if r.HTTPStatusCode != http.StatusCreated {
			t.Error("A request failed")
		}
	}
}

type Application struct {
	LoanAmount         float64 `force:"Requested_Loan_Amount__c"`
	ApplicationNumber  string  `force:"Application_No__c"`
	ClearmatchCustomer string  `force:"Clearmatch_Customer__c"`
	UselessField       int     `force:"useless,omitempty"`
}
