package force

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestComposite_RequestAndResponse(t *testing.T) {
	forceAPI := createTest()

	url := forceAPI.apiSObjects["Folder"]
	urls := url.URLs
	fmt.Println(urls[rowTemplateKey])

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

func TestComposite_Query(t *testing.T) {
	forceAPI := createTest()

	contactID := "0032O000002KnrEQAS" //Replace this with a contact ID that exists :D
	query := fmt.Sprintf("SELECT Id FROM Contact WHERE Id = '%v'", contactID)

	newCompositeRequests := &CompositeRequests{}
	newCompositeRequests.AllOrNone = true

	newRequest := forceAPI.CompositeQuery(query, "Query")

	newCompositeRequests.Add(newRequest)
	resp, err := forceAPI.PostCompositeRequests(newCompositeRequests)

	st, _ := json.Marshal(resp)
	fmt.Println(string(st))

	if err != nil {
		t.Error("Invalid Composite Requests")
	}

	for _, r := range resp.CompositeResponse {
		if r.HTTPStatusCode != http.StatusOK {
			t.Error("A request failed")
		}
	}
}

func TestComposite_Get(t *testing.T) {
	forceAPI := createTest()

	applicationID := "a0u2O0000017azQQAQ" //Replace this with an ID that exists :D

	newCompositeRequests := &CompositeRequests{}
	newCompositeRequests.AllOrNone = true

	newRequest := forceAPI.CompositeGetSObject(applicationID, &Application{}, []string{"Id"}, "Get")

	newCompositeRequests.Add(newRequest)
	resp, err := forceAPI.PostCompositeRequests(newCompositeRequests)

	st, _ := json.Marshal(resp)
	fmt.Println(string(st))

	if err != nil {
		t.Error("Invalid Composite Requests")
	}

	for _, r := range resp.CompositeResponse {
		if r.HTTPStatusCode != http.StatusOK {
			t.Error("A request failed")
		}
	}
}

func TestComposite_Insert(t *testing.T) {
	forceAPI := createTest()

	contactID := "0032O000002KnrEQAS" //Replace this with a contact ID that exists :D

	newCompositeRequests := &CompositeRequests{}
	newCompositeRequests.AllOrNone = true

	newApplication := &Application{
		ClearmatchCustomer: contactID,
		ApplicationNumber:  "A1111113",
		LoanAmount:         1111.11,
	}

	newRequest := forceAPI.CompositeInsertSObject(newApplication, "Insert Application")

	newCompositeRequests.Add(newRequest)
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

func TestComposite_Update(t *testing.T) {
	forceAPI := createTest()

	//Replace this with an ID that exists :D
	applicationID := "a0u2O0000017azQQAQ"

	newCompositeRequests := &CompositeRequests{}
	newCompositeRequests.AllOrNone = true

	newApplication := &Application{
		ApplicationNumber: "A1111113",
		LoanAmount:        1111.11,
	}

	newRequest := forceAPI.CompositeUpdateSObject(applicationID, newApplication, "Update Application")

	newCompositeRequests.Add(newRequest)
	resp, err := forceAPI.PostCompositeRequests(newCompositeRequests)

	st, _ := json.Marshal(resp)
	fmt.Println(string(st))

	if err != nil {
		t.Error("Invalid Composite Requests")
	}

	for _, r := range resp.CompositeResponse {
		if r.HTTPStatusCode != http.StatusNoContent {
			t.Error("A request failed")
		}
	}
}

func TestComposite_Delete(t *testing.T) {
	forceAPI := createTest()

	applicationID := "a0u2O0000017b18QAA" //Replace this with an ID that exists :D

	newCompositeRequests := &CompositeRequests{}
	newCompositeRequests.AllOrNone = true

	newRequest := forceAPI.CompositeDeleteSObject(applicationID, &Application{}, "Delete Application")

	newCompositeRequests.Add(newRequest)
	resp, err := forceAPI.PostCompositeRequests(newCompositeRequests)

	st, _ := json.Marshal(resp)
	fmt.Println(string(st))

	if err != nil {
		t.Error("Invalid Composite Requests")
	}

	for _, r := range resp.CompositeResponse {
		if r.HTTPStatusCode != http.StatusNoContent {
			t.Error("A request failed")
		}
	}
}

type Application struct {
	LoanAmount         float64 `force:"Requested_Loan_Amount__c,omitempty"`
	ApplicationNumber  string  `force:"Application_No__c,omitempty"`
	ClearmatchCustomer string  `force:"Clearmatch_Customer__c,omitempty"`
	UselessField       int     `force:"useless,omitempty"`
}

//ApiName to fulfill the SObject interface
func (q *Application) ApiName() string {
	return "Applications__c"
}

//ExternalIdApiName to fulfill the SObject interface
func (q *Application) ExternalIdApiName() string {
	return "Applications__c"
}
