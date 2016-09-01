package force

import (
	"math/rand"
	"testing"
	"time"

	"github.com/nimajalali/go-force/sobjects"
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

func TestDescribeSobjects(t *testing.T) {
	forceAPI := createTest()
	objects, err := forceAPI.DescribeSObjects()
	if err != nil {
		t.Fatal("Failed to retrieve SObjects", err)
	}
	t.Logf("SObjects for Account Retrieved: %+v", objects)
}

func TestDescribeSObject(t *testing.T) {
	forceApi := createTest()
	acc := &sobjects.Account{}

	desc, err := forceApi.DescribeSObject(acc)
	if err != nil {
		t.Fatalf("Cannot retrieve SObject Description for Account SObject: %v", err)
	}

	t.Logf("SObject Description for Account Retrieved: %+v", desc)
}

func TestGetSObject(t *testing.T) {
	forceApi := createTest()
	// Test Standard Object
	acc := &sobjects.Account{}

	err := forceApi.GetSObject(AccountId, nil, acc)
	if err != nil {
		t.Fatalf("Cannot retrieve SObject Account: %v", err)
	}

	t.Logf("SObject Account Retrieved: %+v", acc)

	// Test Custom Object
	customObject := &CustomSObject{}

	err = forceApi.GetSObject(CustomObjectId, nil, customObject)
	if err != nil {
		t.Fatalf("Cannot retrieve SObject CustomObject: %v", err)
	}

	t.Logf("SObject CustomObject Retrieved: %+v", customObject)

	// Test Custom Object Field Retrieval
	fields := []string{"Name", "Id"}

	accFields := &sobjects.Account{}

	err = forceApi.GetSObject(AccountId, fields, accFields)
	if err != nil {
		t.Fatalf("Cannot retrieve SObject Account fields: %v", err)
	}

	t.Logf("SObject Account Name and Id Retrieved: %+v", accFields)
}

func TestUpdateSObject(t *testing.T) {
	forceApi := createTest()
	// Need some random text for updating a field.
	rand.Seed(time.Now().UTC().UnixNano())
	someText := randomString(10)

	// Test Standard Object
	acc := &sobjects.Account{}
	acc.Name = someText

	err := forceApi.UpdateSObject(AccountId, acc)
	if err != nil {
		t.Fatalf("Cannot update SObject Account: %v", err)
	}

	// Read back and verify
	err = forceApi.GetSObject(AccountId, nil, acc)
	if err != nil {
		t.Fatalf("Cannot retrieve SObject Account: %v", err)
	}

	if acc.Name != someText {
		t.Fatalf("Update SObject Account failed. Failed to persist.")
	}

	t.Logf("Updated SObject Account: %+v", acc)
}

func TestInsertDeleteSObject(t *testing.T) {
	forceApi := createTest()
	objectId := insertSObject(forceApi, t)
	deleteSObject(forceApi, t, objectId)
}

func insertSObject(forceApi *ForceApi, t *testing.T) string {
	// Need some random text for name field.
	rand.Seed(time.Now().UTC().UnixNano())
	someText := randomString(10)

	// Test Standard Object
	acc := &sobjects.Account{}
	acc.Name = someText

	resp, err := forceApi.InsertSObject(acc)
	if err != nil {
		t.Fatalf("Insert SObject Account failed: %v", err)
	}

	if len(resp.Id) == 0 {
		t.Fatalf("Insert SObject Account failed to return Id: %+v", resp)
	}

	return resp.Id
}

func deleteSObject(forceApi *ForceApi, t *testing.T, id string) {
	// Test Standard Object
	acc := &sobjects.Account{}

	err := forceApi.DeleteSObject(id, acc)
	if err != nil {
		t.Fatalf("Delete SObject Account failed: %v", err)
	}

	// Read back and verify
	err = forceApi.GetSObject(id, nil, acc)
	if err == nil {
		t.Fatalf("Delete SObject Account failed, was able to retrieve deleted object: %+v", acc)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
