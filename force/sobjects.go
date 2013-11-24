package force

import (
	"bytes"
	"fmt"
	"strings"
)

// Interface all standard and custom objects must implement. Needed for uri generation.
type SObject interface {
	ApiName() string
	ExternalIdApiName() string
}

// Response recieved from force.com API after insert of an sobject.
type SObjectResponse struct {
	Id      string    `force:"id,omitempty"`
	Errors  ApiErrors `force:"error,omitempty"` //TODO: Not sure if ApiErrors is the right object
	Success bool      `force:"success,omitempty"`
}

func DescribeSObject(in SObject) (resp *SObjectDescription, err error) {
	// Check cache
	resp, ok := apiSObjectDescriptions[in.ApiName()]
	if !ok {
		// Attempt retrieval from api
		sObjectMetaData, ok := apiSObjects[in.ApiName()]
		if !ok {
			err = fmt.Errorf("Unable to find metadata for object: %v", in.ApiName())
			return
		}

		uri := sObjectMetaData.URLs[sObjectDescribeKey]

		resp = &SObjectDescription{}
		err = get(uri, nil, resp)
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

		apiSObjectDescriptions[in.ApiName()] = resp
	}

	return
}

// TODO: Add fields parameter to only retireve needed fields.
func GetSObject(id string, out SObject) (err error) {
	uri := strings.Replace(apiSObjects[out.ApiName()].URLs[rowTemplateKey], idKey, id, 1)

	err = get(uri, nil, out.(interface{}))

	return
}

func InsertSObject(in SObject) (resp *SObjectResponse, err error) {
	uri := apiSObjects[in.ApiName()].URLs[sObjectKey]

	resp = &SObjectResponse{}
	err = post(uri, nil, in.(interface{}), resp)

	return
}

func UpdateSObject(id string, in SObject) (err error) {
	uri := strings.Replace(apiSObjects[in.ApiName()].URLs[rowTemplateKey], idKey, id, 1)

	err = patch(uri, nil, in.(interface{}), nil)

	return
}

func DeleteSObject(id string, in SObject) (err error) {
	uri := strings.Replace(apiSObjects[in.ApiName()].URLs[rowTemplateKey], idKey, id, 1)

	err = delete(uri, nil)

	return
}

// TODO: Add fields parameter to only retireve needed fields.
func GetSObjectByExternalId(id string, out SObject) (err error) {
	uri := fmt.Sprintf("%v/%v/%v", apiSObjects[out.ApiName()].URLs[sObjectKey], out.ExternalIdApiName(), id)

	err = get(uri, nil, out.(interface{}))

	return
}

func UpsertSObjectByExternalId(id string, in SObject) (resp *SObjectResponse, err error) {
	uri := fmt.Sprintf("%v/%v/%v", apiSObjects[in.ApiName()].URLs[sObjectKey], in.ExternalIdApiName(), id)

	resp = &SObjectResponse{}
	err = patch(uri, nil, in.(interface{}), resp)

	return
}

func DeleteSObjectByExternalId(id string, in SObject) (err error) {
	uri := fmt.Sprintf("%v/%v/%v", apiSObjects[in.ApiName()].URLs[sObjectKey], in.ExternalIdApiName(), id)

	err = delete(uri, nil)

	return
}
