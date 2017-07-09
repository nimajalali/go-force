package force

type Limits map[string]Limit

type Limit struct {
	Remaining float64
	Max       float64
}

func (forceAPI *ForceAPI) GetLimits() (limits *Limits, err error) {
	uri := forceAPI.apiResources[limitsKey]

	limits = &Limits{}
	err = forceAPI.Get(uri, nil, limits)

	return
}
