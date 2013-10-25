package force

import (
	"fmt"
	"testing"

	"github.com/nimajalali/go-force/sobjects"
)

const (
	query    = "SELECT Id, Name FROM Account LIMIT 10"
	queryAll = "SELECT Id, Name FROM Account WHERE Id = '%v'"
)

func init() {
	initTest()
}

type AccountQueryResponse struct {
	BaseQuery
	Records []sobjects.Account `json:"records" force:"records"`
}

func TestQuery(t *testing.T) {
	list := &AccountQueryResponse{}
	err := Query(query, list)
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
	list := &AccountQueryResponse{}
	err := QueryAll(fmt.Sprintf(queryAll, newId), list)
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
