package force

// Limits is map containing limits.
type Limits map[string]Limit

// Limit describes an API limit.
type Limit struct {
	Remaining float64
	Max       float64
}

// GetLimits returns a specific API limit.
func (forceAPI *API) GetLimits() (limits *Limits, err error) {
	uri := forceAPI.apiResources[limitsKey]

	limits = &Limits{}
	err = forceAPI.Get(uri, nil, limits)

	return
}
