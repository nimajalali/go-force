package force

type Limits map[string]Limit

type Limit struct {
	Remaining float64
	Max       float64
}

func (forceApi *ForceApi) GetLimits() (limits *Limits, err error) {
	uri := forceApi.apiResources[limitsKey]

	limits = &Limits{}
	err = forceApi.Get(uri, nil, limits)

	return
}
