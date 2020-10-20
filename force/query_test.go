package force

import (
	"context"
	"fmt"
	"testing"

	"github.com/EverlongProject/go-force/sobjects"
)

const (
	queryAll = "SELECT %v FROM Account WHERE Id = '%v'"
)

type AccountQueryResponse struct {
	sobjects.BaseQuery
	Records []sobjects.Account `json:"Records" force:"records"`
}

func TestQuery(t *testing.T) {
	forceApi := createTest()
	ctx := context.Background()
	desc, err := forceApi.DescribeSObject(ctx, &sobjects.Account{})
	if err != nil {
		t.Fatalf("Failed to retrieve description of sobject: %v", err)
	}

	list := &AccountQueryResponse{}
	err = forceApi.Query(ctx, BuildQuery(desc.AllFields, desc.Name, nil), list)
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	t.Logf("%#v", list)
}

func TestQueryAll(t *testing.T) {
	forceApi := createTest()
	ctx := context.Background()
	// First Insert and Delete an Account
	newId := insertSObject(ctx, forceApi, t)
	deleteSObject(ctx, forceApi, t, newId)

	// Then look for it.
	desc, err := forceApi.DescribeSObject(ctx, &sobjects.Account{})
	if err != nil {
		t.Fatalf("Failed to retrieve description of sobject: %v", err)
	}

	list := &AccountQueryResponse{}
	err = forceApi.QueryAll(ctx, fmt.Sprintf(queryAll, desc.AllFields, newId), list)
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
