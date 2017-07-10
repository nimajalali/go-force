package force

import (
	"fmt"
	"testing"

	"github.com/nimajalali/go-force/sobjects"
)

const (
	queryAll = "SELECT %v FROM Account WHERE Id = '%v'"
)

type AccountQueryResponse struct {
	sobjects.BaseQuery
	Records []sobjects.Account `json:"Records" force:"records"`
}

func TestQuery(t *testing.T) {
	forceAPI := createTest()
	desc, err := forceAPI.DescribeSObject(&sobjects.Account{})
	if err != nil {
		t.Fatalf("Failed to retrieve description of sobject: %v", err)
	}

	list := &AccountQueryResponse{}
	err = forceAPI.Query(BuildQuery(desc.AllFields, desc.Name, nil), list)
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	t.Logf("%#v", list)
}

func TestQueryAll(t *testing.T) {
	forceAPI := createTest()
	// First Insert and Delete an Account
	newId := insertSObject(forceAPI, t)
	deleteSObject(forceAPI, t, newId)

	// Then look for it.
	desc, err := forceAPI.DescribeSObject(&sobjects.Account{})
	if err != nil {
		t.Fatalf("Failed to retrieve description of sobject: %v", err)
	}

	list := &AccountQueryResponse{}
	err = forceAPI.QueryAll(fmt.Sprintf(queryAll, desc.AllFields, newId), list)
	if err != nil {
		t.Fatalf("Failed to queryAll: %v", err)
	}

	if len(list.Records) == 0 {
		t.Fatal("Failed to retrieve deleted record using queryAll")
	}

	t.Logf("%#v", list)
}

func TestQueryNext(t *testing.T) {
	// TODO
}
