package force

import (
	"testing"

	"go-force/sobjects"
)

const (
	AccountId      = "001i000000RxW18"
	CustomObjectId = "a00i0000009SPer"
)

type CustomSObject struct {
	sobjects.BaseSObject
	Active    bool   `force:"Active__c"`
	AccountId string `force:"Account__c"`
}

func (t *CustomSObject) ApiName() string {
	return "CustomObject__c"
}

func TestGetSObject(t *testing.T) {
	// Test Standard Object
	acc := &sobjects.Account{}

	err := GetSObject(AccountId, acc)
	if err != nil {
		t.Fatalf("Cannot retrieve SObject Account: %v", err)
	}

	t.Logf("SObject Account Retrieved: %#v", acc)

	// Test Custom Object
	customObject := &CustomSObject{}

	err = GetSObject(CustomObjectId, customObject)
	if err != nil {
		t.Fatalf("Cannot retrieve SObject CustomObject: %v", err)
	}

	t.Logf("SObject CustomObject Retrieved: %#v", customObject)
}

func TestUpdateSObject(t *testing.T) {

	// Test Standard Object
	acc := &sobj
}

func insertSObject(t *testing.T) string {

}

func deleteSObject(t *testing.T, id string) {

}
