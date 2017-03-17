package force

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

// SObject is an interface that all standard and custom objects must implement. Needed for uri generation.
type SObject interface {
	APIName() string
	ExternalIDAPIName() string
}

// SObjectResponse represents a response received from force.com API after insert of an SObject.
type SObjectResponse struct {
	ID      string    `force:"id,omitempty"`
	Errors  APIErrors `force:"error,omitempty"` //TODO: Not sure if APIErrors is the right object
	Success bool      `force:"success,omitempty"`
}

// DescribeSObjects describes a list of SObjects.
func (forceAPI *API) DescribeSObjects() (map[string]*SObjectMetaData, error) {
	if err := forceAPI.getAPISObjects(); err != nil {
		return nil, err
	}

	return forceAPI.apiSObjects, nil
}

// DescribeSObject describes an SObject.
func (forceAPI *API) DescribeSObject(in SObject) (resp *SObjectDescription, err error) {
	// Check cache
	resp, ok := forceAPI.apiSObjectDescriptions[in.APIName()]
	if !ok {
		// Attempt retrieval from api
		sObjectMetaData, ok := forceAPI.apiSObjects[in.APIName()]
		if !ok {
			err = fmt.Errorf("Unable to find metadata for object: %v", in.APIName())
			return
		}

		uri := sObjectMetaData.URLs[sObjectDescribeKey]

		resp = &SObjectDescription{}
		err = forceAPI.Get(uri, nil, resp)
		if err != nil {
			return
		}

		// Create Comma Separated String of All Field Names.
		// Used for SELECT * Queries.
		length := len(resp.Fields)
		if length > 0 {
			var allFields bytes.Buffer
			for index, field := range resp.Fields {
				// Field type location cannot be directly retrieved from SQL Query.
				if field.Type != "location" {
					if index > 0 && index < length {
						allFields.WriteString(", ")
					}
					allFields.WriteString(field.Name)
				}
			}

			resp.AllFields = allFields.String()
		}

		forceAPI.apiSObjectDescriptions[in.APIName()] = resp
	}

	return
}

// GetSObject returns an SObject.
func (forceAPI *API) GetSObject(id string, fields []string, out SObject) (err error) {
	uri := strings.Replace(forceAPI.apiSObjects[out.APIName()].URLs[rowTemplateKey], idKey, id, 1)

	params := url.Values{}
	if len(fields) > 0 {
		params.Add("fields", strings.Join(fields, ","))
	}

	err = forceAPI.Get(uri, params, out.(interface{}))

	return
}

// InsertSObject creates an SObject.
func (forceAPI *API) InsertSObject(in SObject) (resp *SObjectResponse, err error) {
	uri := forceAPI.apiSObjects[in.APIName()].URLs[sObjectKey]

	resp = &SObjectResponse{}
	err = forceAPI.Post(uri, nil, in.(interface{}), resp)

	return
}

// UpdateSObject updates an SObject.
func (forceAPI *API) UpdateSObject(id string, in SObject) (err error) {
	uri := strings.Replace(forceAPI.apiSObjects[in.APIName()].URLs[rowTemplateKey], idKey, id, 1)

	err = forceAPI.Patch(uri, nil, in.(interface{}), nil)

	return
}

// DeleteSObject deletes an SObject.
func (forceAPI *API) DeleteSObject(id string, in SObject) (err error) {
	uri := strings.Replace(forceAPI.apiSObjects[in.APIName()].URLs[rowTemplateKey], idKey, id, 1)

	err = forceAPI.Delete(uri, nil)

	return
}

// GetSObjectByExternalID returns an SObject by external ID.
func (forceAPI *API) GetSObjectByExternalID(id string, fields []string, out SObject) (err error) {
	uri := fmt.Sprintf("%v/%v/%v", forceAPI.apiSObjects[out.APIName()].URLs[sObjectKey],
		out.ExternalIDAPIName(), id)

	params := url.Values{}
	if len(fields) > 0 {
		params.Add("fields", strings.Join(fields, ","))
	}

	err = forceAPI.Get(uri, params, out.(interface{}))

	return
}

// UpsertSObjectByExternalID performs an upsert using an external ID.
func (forceAPI *API) UpsertSObjectByExternalID(id string, in SObject) (resp *SObjectResponse, err error) {

	uri := fmt.Sprintf("%v/%v/%v", forceAPI.apiSObjects[in.APIName()].URLs[sObjectKey],
		in.ExternalIDAPIName(), id)

	resp = &SObjectResponse{}
	err = forceAPI.Patch(uri, nil, in.(interface{}), resp)

	return
}

// UpsertSObjectStringByExternalID performs an upsert using an
// external ID and the JSON string representation of an SObject.
func (forceAPI *API) UpsertSObjectStringByExternalID(object, extenalID, id, data string) (resp *SObjectResponse, err error) {

	uri := fmt.Sprintf("%v/%v/%v", forceAPI.apiSObjects[object].URLs[sObjectKey], extenalID, id)

	resp = &SObjectResponse{}
	err = forceAPI.Patch(uri, nil, data, resp)

	return
}

// DeleteSObjectByExternalID deletes an SObject by external ID.
func (forceAPI *API) DeleteSObjectByExternalID(id string, in SObject) (err error) {
	uri := fmt.Sprintf("%v/%v/%v", forceAPI.apiSObjects[in.APIName()].URLs[sObjectKey],
		in.ExternalIDAPIName(), id)

	err = forceAPI.Delete(uri, nil)

	return
}
