package force

import "net/url"

// Search uses the Search resource to execute an SOSL query and write all the results to
// the `out` argument
func (forceAPI *ForceAPI) Search(query string, out interface{}) error {
	uri := forceAPI.apiResources[searchKey]

	params := url.Values{
		"q": {query},
	}

	return forceAPI.Get(uri, params, out)
}
