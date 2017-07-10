package force

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	BaseQueryString = "SELECT %v FROM %v"
)

func BuildQuery(fields, table string, constraints []string) string {
	query := fmt.Sprintf(BaseQueryString, fields, table)
	if len(constraints) > 0 {
		query += fmt.Sprintf(" WHERE %v", strings.Join(constraints, " AND "))
	}

	return query
}

// Query uses the Query resource to execute an SOQL query that returns all the results in a single response,
// or if needed, returns part of the results and an identifier used to retrieve the remaining results.
func (forceAPI *ForceAPI) Query(query string, out interface{}) error {
	uri := forceAPI.apiResources[queryKey]

	params := url.Values{
		"q": {query},
	}

	return forceAPI.Get(uri, params, out)
}

// Use the QueryAll resource to execute a SOQL query that includes information about records that have
// been deleted because of a merge or delete. Use QueryAll rather than Query, because the Query resource
// will automatically filter out items that have been deleted.
func (forceAPI *ForceAPI) QueryAll(query string, out interface{}) (err error) {
	uri := forceAPI.apiResources[queryAllKey]

	params := url.Values{
		"q": {query},
	}

	err = forceAPI.Get(uri, params, out)

	return
}

func (forceAPI *ForceAPI) QueryNext(uri string, out interface{}) (err error) {
	err = forceAPI.Get(uri, nil, out)

	return
}
