package force

import (
	"fmt"
)

const (
	limitsUri = "/services/data/%v/limits"
)

type Limits map[string]Limit

type Limit struct {
	Remaining float64
	Max       float64
}

func GetLimits() (limits *Limits, err error) {
	uri := fmt.Sprintf(limitsUri, apiVersion)

	limits = &Limits{}
	err = get(uri, nil, limits)

	return
}
