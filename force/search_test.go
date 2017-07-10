package force

import (
	"testing"

	"github.com/goguardian/go-force/sobjects"
)

type AccountsSearchResponse []sobjects.Account

func TestSearch(t *testing.T) {
	forceAPI := createTest()

	list := &AccountsSearchResponse{}
	err := forceAPI.Search("FIND {all} IN ALL FIELDS RETURNING Account (Id)", list)
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	t.Logf("%#v", list)
}
