package force

import (
	"fmt"
	"testing"

	"github.com/nimajalali/go-force/sobjects"
)

const (
	queryAll = "SELECT %v FROM Account WHERE Id = '%v'"
)

func init() {
	initTest()
}

type AccountQueryResponse struct {
	sobjects.BaseQuery
	Records []sobjects.Account `json:"Records" force:"records"`
}

func TestQuery(t *testing.T) {
	desc, err := DescribeSObject(&sobjects.Account{})
	if err != nil {
		t.Fatalf("Failed to retrieve description of sobject: %v", err)
	}

	list := &AccountQueryResponse{}
	err = Query(fmt.Sprintf(BaseQueryString, desc.AllFields, desc.Name), list)
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	t.Logf("%#v", list)
}

func TestQueryAll(t *testing.T) {
	// First Insert and Delete an Account
	newId := insertSObject(t)
	deleteSObject(t, newId)

	// Then look for it.
	desc, err := DescribeSObject(&sobjects.Account{})
	if err != nil {
		t.Fatalf("Failed to retrieve description of sobject: %v", err)
	}

	list := &AccountQueryResponse{}
	err = QueryAll(fmt.Sprintf(queryAll, desc.AllFields, newId), list)
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
