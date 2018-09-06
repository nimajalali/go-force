package force

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/nimajalali/go-force/forcejson"
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

func TestUnmarshalCompositResponse(t *testing.T) {
	response := `{
		"compositeResponse": [
			{
				"body": {
					"id": "a0u2O0000017bNpQAI",
					"success": true,
					"errors": []
				},
				"httpHeaders": {
					"Location": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017bNpQAI"
				},
				"httpStatusCode": 201,
				"referenceId": "NewApp"
			},
			{
				"body": {
					"id": "a0u2O0000017bNqQAI",
					"success": true,
					"errors": []
				},
				"httpHeaders": {
					"Location": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017bNqQAI"
				},
				"httpStatusCode": 201,
				"referenceId": "NewApp1"
			},
			{
				"body": {
					"totalSize": 30,
					"done": true,
					"records": [
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017azQQAQ"
							},
							"Id": "a0u2O0000017azQQAQ"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017azVQAQ"
							},
							"Id": "a0u2O0000017azVQAQ"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017azaQAA"
							},
							"Id": "a0u2O0000017azaQAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017azfQAA"
							},
							"Id": "a0u2O0000017azfQAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017azpQAA"
							},
							"Id": "a0u2O0000017azpQAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017azqQAA"
							},
							"Id": "a0u2O0000017azqQAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017azuQAA"
							},
							"Id": "a0u2O0000017azuQAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017azzQAA"
							},
							"Id": "a0u2O0000017azzQAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0EQAQ"
							},
							"Id": "a0u2O0000017b0EQAQ"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0FQAQ"
							},
							"Id": "a0u2O0000017b0FQAQ"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0JQAQ"
							},
							"Id": "a0u2O0000017b0JQAQ"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0KQAQ"
							},
							"Id": "a0u2O0000017b0KQAQ"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0OQAQ"
							},
							"Id": "a0u2O0000017b0OQAQ"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0PQAQ"
							},
							"Id": "a0u2O0000017b0PQAQ"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0TQAQ"
							},
							"Id": "a0u2O0000017b0TQAQ"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0UQAQ"
							},
							"Id": "a0u2O0000017b0UQAQ"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0YQAQ"
							},
							"Id": "a0u2O0000017b0YQAQ"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0ZQAQ"
							},
							"Id": "a0u2O0000017b0ZQAQ"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0dQAA"
							},
							"Id": "a0u2O0000017b0dQAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0eQAA"
							},
							"Id": "a0u2O0000017b0eQAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0nQAA"
							},
							"Id": "a0u2O0000017b0nQAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0oQAA"
							},
							"Id": "a0u2O0000017b0oQAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0sQAA"
							},
							"Id": "a0u2O0000017b0sQAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0tQAA"
							},
							"Id": "a0u2O0000017b0tQAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0xQAA"
							},
							"Id": "a0u2O0000017b0xQAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b0yQAA"
							},
							"Id": "a0u2O0000017b0yQAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017b17QAA"
							},
							"Id": "a0u2O0000017b17QAA"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017bDLQAY"
							},
							"Id": "a0u2O0000017bDLQAY"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017bNpQAI"
							},
							"Id": "a0u2O0000017bNpQAI"
						},
						{
							"attributes": {
								"type": "Applications__c",
								"url": "/services/data/v42.0/sobjects/Applications__c/a0u2O0000017bNqQAI"
							},
							"Id": "a0u2O0000017bNqQAI"
						}
					]
				},
				"httpHeaders": {},
				"httpStatusCode": 200,
				"referenceId": "GetApps"
			}
		]
	}`

	buf := bytes.NewBufferString(response)
	out := &CompositeResponses{}
	if err := forcejson.Unmarshal(buf.Bytes(), out); err != nil {
		t.Error(err)
	}
	expected := 3
	if len(out.CompositeResponse) < expected {
		t.Errorf("Expected CompositeResponses to contain:%v CompositeResponse objects, but got:%v", expected, len(out.CompositeResponse))
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
