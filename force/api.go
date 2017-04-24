package force

import (
	"fmt"
)

const (
	limitsKey          = "limits"
	queryKey           = "query"
	queryAllKey        = "queryAll"
	sObjectsKey        = "sobjects"
	sObjectKey         = "sobject"
	sObjectDescribeKey = "describe"

	rowTemplateKey = "rowTemplate"
	idKey          = "{ID}"

	resourcesURI = "/services/data/%v"
)

// An API represents a SalesForce API client.
type API struct {
	apiVersion             string
	oauth                  *forceOauth
	apiResources           map[string]string
	apiSObjects            map[string]*SObjectMetaData
	apiSObjectDescriptions map[string]*SObjectDescription
	apiMaxBatchSize        int64
	logger                 APILogger
	logPrefix              string
	// OpenJobs contains a list of open jobs mapped by ID.
	OpenJobs map[string]*SJob
}

// RefreshTokenResponse represents a refresh token.
type RefreshTokenResponse struct {
	ID          string `json:"id"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`
	AccessToken string `json:"access_token"`
}

// SObjectAPIResponse represents an API response.
type SObjectAPIResponse struct {
	Encoding     string             `json:"encoding"`
	MaxBatchSize int64              `json:"maxBatchSize"`
	SObjects     []*SObjectMetaData `json:"sobjects"`
}

// SObjectMetaData represents a list of metadata for an SObject.
type SObjectMetaData struct {
	Name                string            `json:"name"`
	Label               string            `json:"label"`
	KeyPrefix           string            `json:"keyPrefix"`
	LabelPlural         string            `json:"labelPlural"`
	Custom              bool              `json:"custom"`
	Layoutable          bool              `json:"layoutable"`
	Activateable        bool              `json:"activateable"`
	URLs                map[string]string `json:"urls"`
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

// SObjectDescription represents an SObject description.
type SObjectDescription struct {
	Name                string               `json:"name"`
	Fields              []*SObjectField      `json:"fields"`
	KeyPrefix           string               `json:"keyPrefix"`
	Layoutable          bool                 `json:"layoutable"`
	Activateable        bool                 `json:"activateable"`
	LabelPlural         string               `json:"labelPlural"`
	Custom              bool                 `json:"custom"`
	CompactLayoutable   bool                 `json:"compactLayoutable"`
	Label               string               `json:"label"`
	Searchable          bool                 `json:"searchable"`
	URLs                map[string]string    `json:"urls"`
	Queryable           bool                 `json:"queryable"`
	Deletable           bool                 `json:"deletable"`
	Updateable          bool                 `json:"updateable"`
	Createable          bool                 `json:"createable"`
	CustomSetting       bool                 `json:"customSetting"`
	Undeletable         bool                 `json:"undeletable"`
	Mergeable           bool                 `json:"mergeable"`
	Replicateable       bool                 `json:"replicateable"`
	Triggerable         bool                 `json:"triggerable"`
	FeedEnabled         bool                 `json:"feedEnabled"`
	Retrievable         bool                 `json:"retrieveable"`
	SearchLayoutable    bool                 `json:"searchLayoutable"`
	LookupLayoutable    bool                 `json:"lookupLayoutable"`
	Listviewable        bool                 `json:"listviewable"`
	DeprecatedAndHidden bool                 `json:"deprecatedAndHidden"`
	RecordTypeInfos     []*RecordTypeInfo    `json:"recordTypeInfos"`
	ChildRelationsips   []*ChildRelationship `json:"childRelationships"`

	AllFields string `json:"-"` // Not from force.com API. Used to generate SELECT * queries.
}

// SObjectField represents a field in an SObject.
type SObjectField struct {
	Length                   float64          `json:"length"`
	Name                     string           `json:"name"`
	Type                     string           `json:"type"`
	RestrictedPicklist       bool             `json:"restrictedPicklist"`
	NameField                bool             `json:"nameField"`
	ByteLength               float64          `json:"byteLength"`
	Precision                float64          `json:"precision"`
	Filterable               bool             `json:"filterable"`
	Sortable                 bool             `json:"sortable"`
	Unique                   bool             `json:"unique"`
	CaseSensitive            bool             `json:"caseSensitive"`
	Calculated               bool             `json:"calculated"`
	Scale                    float64          `json:"scale"`
	Label                    string           `json:"label"`
	NamePointing             bool             `json:"namePointing"`
	Custom                   bool             `json:"custom"`
	HTMLFormatted            bool             `json:"htmlFormatted"`
	DependentPicklist        bool             `json:"dependentPicklist"`
	Permissionable           bool             `json:"permissionable"`
	ReferenceTo              []string         `json:"referenceTo"`
	RelationshipOrder        float64          `json:"relationshipOrder"`
	SoapType                 string           `json:"soapType"`
	CalculatedValueFormula   string           `json:"calculatedValueFormula"`
	DefaultValueFormula      string           `json:"defaultValueFormula"`
	DefaultedOnCreate        bool             `json:"defaultedOnCreate"`
	Digits                   float64          `json:"digits"`
	Groupable                bool             `json:"groupable"`
	Nillable                 bool             `json:"nillable"`
	InlineHelpText           string           `json:"inlineHelpText"`
	WriteRequiresMasterRead  bool             `json:"writeRequiresMasterRead"`
	PicklistValues           []*PicklistValue `json:"picklistValues"`
	Updateable               bool             `json:"updateable"`
	Createable               bool             `json:"createable"`
	DeprecatedAndHidden      bool             `json:"deprecatedAndHidden"`
	DisplayLocationInDecimal bool             `json:"displayLocationInDecimal"`
	CascadeDelete            bool             `json:"cascasdeDelete"`
	RestrictedDelete         bool             `json:"restrictedDelete"`
	ControllerName           string           `json:"controllerName"`
	ExternalID               bool             `json:"externalID"`
	IDLookup                 bool             `json:"idLookup"`
	AutoNumber               bool             `json:"autoNumber"`
	RelationshipName         string           `json:"relationshipName"`
}

// PicklistValue represents a picklist value type.
type PicklistValue struct {
	Value       string `json:"value"`
	DefaulValue bool   `json:"defaultValue"`
	ValidFor    string `json:"validFor"`
	Active      bool   `json:"active"`
	Label       string `json:"label"`
}

// RecordTypeInfo represents a record value type.
type RecordTypeInfo struct {
	Name                     string            `json:"name"`
	Available                bool              `json:"available"`
	RecordTypeID             string            `json:"recordTypeID"`
	URLs                     map[string]string `json:"urls"`
	DefaultRecordTypeMapping bool              `json:"defaultRecordTypeMapping"`
}

// ChildRelationship represents a child relationship for two SObjects.
type ChildRelationship struct {
	Field               string `json:"field"`
	ChildSObject        string `json:"childSObject"`
	DeprecatedAndHidden bool   `json:"deprecatedAndHidden"`
	CascadeDelete       bool   `json:"cascadeDelete"`
	RestrictedDelete    bool   `json:"restrictedDelete"`
	RelationshipName    string `json:"relationshipName"`
}

func (forceAPI *API) getAPIResources() error {
	uri := fmt.Sprintf(resourcesURI, forceAPI.apiVersion)

	return forceAPI.Get(uri, nil, &forceAPI.apiResources)
}

func (forceAPI *API) getAPISObjects() error {
	uri := forceAPI.apiResources[sObjectsKey]

	list := &SObjectAPIResponse{}
	err := forceAPI.Get(uri, nil, list)
	if err != nil {
		return err
	}

	forceAPI.apiMaxBatchSize = list.MaxBatchSize

	// The API doesn't return the list of sobjects in a map. Convert it.
	for _, object := range list.SObjects {
		forceAPI.apiSObjects[object.Name] = object
	}

	return nil
}

func (forceAPI *API) getAPISObjectDescriptions() error {
	for name, metaData := range forceAPI.apiSObjects {
		if _, ok := requestedObjMetadata[name]; !ok {
			continue
		}
		uri := metaData.URLs[sObjectDescribeKey]

		desc := &SObjectDescription{}
		err := forceAPI.Get(uri, nil, desc)
		if err != nil {
			return err
		}

		forceAPI.apiSObjectDescriptions[name] = desc
	}

	return nil
}

// GetInstanceURL returns the instance URL.
func (forceAPI *API) GetInstanceURL() string {
	return forceAPI.oauth.InstanceURL
}

// GetAccessToken returns the access token.
func (forceAPI *API) GetAccessToken() string {
	return forceAPI.oauth.AccessToken
}

// RefreshToken refreshes an access token.
func (forceAPI *API) RefreshToken() error {
	res := &RefreshTokenResponse{}
	payload := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": forceAPI.oauth.refreshToken,
		"client_id":     forceAPI.oauth.clientID,
		"client_secret": forceAPI.oauth.clientSecret,
	}

	err := forceAPI.Post("/services/oauth2/token", nil, payload, res)
	if err != nil {
		return err
	}

	forceAPI.oauth.AccessToken = res.AccessToken
	return nil
}
