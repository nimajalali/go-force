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

var ApiResources map[string]string
var ApiSObjects map[string]SObjectMetaData
var ApiMaxBatchSize int64

func getApiResources() error {
	uri := fmt.Sprintf(resourcesUri, apiVersion)

	ApiResources = make(map[string]string)
	return get(uri, nil, &ApiResources)
}

func getApiSObjects() error {
	uri := ApiResources[sObjectsKey]

	list := &sObjectApiResponse{}
	err := get(uri, nil, list)
	if err != nil {
		return err
	}

	ApiMaxBatchSize = list.MaxBatchSize

	// The API doesn't return the list of sobjects in a map. Convert it.
	ApiSObjects = make(map[string]SObjectMetaData)
	for _, object := range list.SObjects {
		ApiSObjects[object.Name] = object
	}

	return nil
}

type sObjectApiResponse struct {
	Encoding     string            `json:"encoding"`
	MaxBatchSize int64             `json:"maxBatchSize"`
	SObjects     []SObjectMetaData `json:"sobjects"`
}

type SObjectMetaData struct {
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

// Custom Error to handle salesforce api responses.
type ApiErrors []ApiError

type ApiError struct {
	Fields           []string `json:"fields,omitempty" force:"fields,omitempty"`
	Message          string   `json:"message,omitempty" force:"message,omitempty"`
	ErrorCode        string   `json:"errorCode,omitempty" force:"errorCode,omitempty"`
	ErrorName        string   `json:"error,omitempty" force:"error,omitempty"`
	ErrorDescription string   `json:"error_description,omitempty" force:"error_description,omitempty"`
}

func (e ApiErrors) Error() string {
	return fmt.Sprintf("%#v", e)
}

func (e ApiError) Error() string {
	return fmt.Sprintf("%#v", e)
}
