package force

import (
	"fmt"
)

const (
	limitsKey   = "limits"
	queryKey    = "query"
	queryAllKey = "queryAll"
	sObjectsKey = "sobjects"
	sObjectKey  = "sobject"

	rowTemplateKey = "rowTemplate"
	idKey          = "{ID}"

	resourcesUri = "/services/data/%v"
)

var apiResources map[string]string
var apiSObjects map[string]sObjectMetaData
var apiMaxBatchSize int64

func getApiResources() error {
	uri := fmt.Sprintf(resourcesUri, apiVersion)

	apiResources = make(map[string]string)
	return get(uri, nil, &apiResources)
}

func getApiSObjects() error {
	uri := apiResources[sObjectsKey]

	list := &sObjectApiResponse{}
	err := get(uri, nil, list)
	if err != nil {
		return err
	}

	apiMaxBatchSize = list.MaxBatchSize

	// The API doesn't return the list of sobjects in a map. Convert it.
	apiSObjects = make(map[string]sObjectMetaData)
	for _, object := range list.SObjects {
		apiSObjects[object.Name] = object
	}

	return nil
}

type sObjectApiResponse struct {
	Encoding     string            `json:"encoding"`
	MaxBatchSize int64             `json:"maxBatchSize"`
	SObjects     []sObjectMetaData `json:"sobjects"`
}

type sObjectMetaData struct {
	Name                string            `json:"name"`
	Label               string            `json:"label"`
	KeyPrefix           string            `json:"keyPrefix"`
	LabelPlural         string            `json:"labelPlural"`
	Custom              bool              `json:"custom"`
	Layoutable          bool              `json:"layoutable"`
	Activateable        bool              `json:"activateable"`
	Urls                map[string]string `json:"urls"`
	Searchable          bool              `json:"searchable"`
	Updateable          bool              `json:"updateable"`
	Createable          bool              `json:"createable"`
	DeprecatedAndHidden bool              `json:"deprecatedAndHidden"`
	CustomSetting       bool              `json:"customSetting"`
	Deletable           bool              `json:"deletable"`
	FeedEnabled         bool              `json:"feedEnabled"`
	Mergeable           bool              `json:"mergeable"`
	Queryable           bool              `json:"queryable"`
	Replicateable       bool              `json:"replicateable"`
	Retrieveable        bool              `json:"retrieveable"`
	Undeletable         bool              `json:"undeletable"`
	Triggerable         bool              `json:"triggerable"`
}
