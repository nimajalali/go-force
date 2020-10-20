package force

import "context"

type Limits map[string]Limit

type Limit struct {
	Remaining float64
	Max       float64
}

func (forceApi *ForceApi) GetLimits(ctx context.Context) (limits *Limits, err error) {
	uri := forceApi.apiResources[limitsKey]

	limits = &Limits{}
	err = forceApi.Get(ctx, uri, nil, limits)

	return
}
