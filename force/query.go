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

// Use the Query resource to execute a SOQL query that returns all the results in a single response,
// or if needed, returns part of the results and an identifier used to retrieve the remaining results.
func (forceApi *ForceApi) Query(query string, out interface{}) (err error) {
	uri := forceApi.apiResources[queryKey]

	params := url.Values{
		"q": {query},
	}

	err = forceApi.Get(uri, params, out)

	return
}

// Use the QueryAll resource to execute a SOQL query that includes information about records that have
// been deleted because of a merge or delete. Use QueryAll rather than Query, because the Query resource
// will automatically filter out items that have been deleted.
func (forceApi *ForceApi) QueryAll(query string, out interface{}) (err error) {
	uri := forceApi.apiResources[queryAllKey]

	params := url.Values{
		"q": {query},
	}

	err = forceApi.Get(uri, params, out)

	return
}

func (forceApi *ForceApi) QueryNext(uri string, out interface{}) (err error) {
	err = forceApi.Get(uri, nil, out)

	return
}
