package force

import (
	"net/url"
)

// Use the Query resource to execute a SOQL query that returns all the results in a single response,
// or if needed, returns part of the results and an identifier used to retrieve the remaining results.
func Query(query string, out interface{}) (err error) {
	uri := apiResources[queryKey]

	params := url.Values{
		"q": {query},
	}

	err = get(uri, params, out)

	return
}

// Use the QueryAll resource to execute a SOQL query that includes information about records that have
// been deleted because of a merge or delete. Use QueryAll rather than Query, because the Query resource
// will automatically filter out items that have been deleted.
func QueryAll(query string, out interface{}) (err error) {
	uri := apiResources[queryAllKey]

	params := url.Values{
		"q": {query},
	}

	err = get(uri, params, out)

	return
}

func QueryNext(uri string, out interface{}) (err error) {
	err = get(uri, nil, out)

	return
}
