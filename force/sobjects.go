package force

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

// Interface all standard and custom objects must implement. Needed for uri generation.
type SObject interface {
	APIName() string
	ExternalIDAPIName() string
}

// Response received from force.com API after insert of an sobject.
type SObjectResponse struct {
	Id      string    `force:"id,omitempty"`
	Errors  APIErrors `force:"error,omitempty"` //TODO: Not sure if APIErrors is the right object
	Success bool      `force:"success,omitempty"`
}

func (forceAPI *ForceAPI) DescribeSObjects() (map[string]*SObjectMetaData, error) {
	if err := forceAPI.getSObjects(); err != nil {
		return nil, err
	}

	return forceAPI.apiSObjects, nil
}

func (forceAPI *ForceAPI) DescribeSObject(in SObject) (resp *SObjectDescription, err error) {
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

func (forceAPI *ForceAPI) GetSObject(id string, fields []string, out SObject) (err error) {
	uri := strings.Replace(forceAPI.apiSObjects[out.APIName()].URLs[rowTemplateKey], idKey, id, 1)

	params := url.Values{}
	if len(fields) > 0 {
		params.Add("fields", strings.Join(fields, ","))
	}

	err = forceAPI.Get(uri, params, out.(interface{}))

	return
}

func (forceAPI *ForceAPI) InsertSObject(in SObject) (resp *SObjectResponse, err error) {
	uri := forceAPI.apiSObjects[in.APIName()].URLs[sObjectKey]

	resp = &SObjectResponse{}
	err = forceAPI.Post(uri, nil, in.(interface{}), resp)

	return
}

func (forceAPI *ForceAPI) UpdateSObject(id string, in SObject) (err error) {
	uri := strings.Replace(forceAPI.apiSObjects[in.APIName()].URLs[rowTemplateKey], idKey, id, 1)

	err = forceAPI.Patch(uri, nil, in.(interface{}), nil)

	return
}

func (forceAPI *ForceAPI) DeleteSObject(id string, in SObject) (err error) {
	uri := strings.Replace(forceAPI.apiSObjects[in.APIName()].URLs[rowTemplateKey], idKey, id, 1)

	err = forceAPI.Delete(uri, nil)

	return
}

func (forceAPI *ForceAPI) GetSObjectByExternalId(id string, fields []string, out SObject) (err error) {
	uri := fmt.Sprintf("%v/%v/%v", forceAPI.apiSObjects[out.APIName()].URLs[sObjectKey],
		out.ExternalIDAPIName(), id)

	params := url.Values{}
	if len(fields) > 0 {
		params.Add("fields", strings.Join(fields, ","))
	}

	err = forceAPI.Get(uri, params, out.(interface{}))

	return
}

func (forceAPI *ForceAPI) UpsertSObjectByExternalId(id string, in SObject) (resp *SObjectResponse, err error) {
	uri := fmt.Sprintf("%v/%v/%v", forceAPI.apiSObjects[in.APIName()].URLs[sObjectKey],
		in.ExternalIDAPIName(), id)

	resp = &SObjectResponse{}
	err = forceAPI.Patch(uri, nil, in.(interface{}), resp)

	return
}

func (forceAPI *ForceAPI) DeleteSObjectByExternalId(id string, in SObject) (err error) {
	uri := fmt.Sprintf("%v/%v/%v", forceAPI.apiSObjects[in.APIName()].URLs[sObjectKey],
		in.ExternalIDAPIName(), id)

	err = forceAPI.Delete(uri, nil)

	return
}
