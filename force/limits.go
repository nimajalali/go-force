package force

type Limits map[string]Limit

type Limit struct {
	Remaining float64
	Max       float64
}

func GetLimits() (limits *Limits, err error) {
	uri := ApiResources[limitsKey]

	limits = &Limits{}
	err = get(uri, nil, limits)

	return
}
