package force

import (
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

// TODO: Add fields parameter to only retireve needed fields.
func GetSObject(id string, out SObject) (err error) {
	uri := strings.Replace(ApiSObjects[out.ApiName()].Urls[rowTemplateKey], idKey, id, 1)

	err = get(uri, nil, out.(interface{}))

	return
}

func InsertSObject(in SObject) (resp *SObjectResponse, err error) {
	uri := ApiSObjects[in.ApiName()].Urls[sObjectKey]

	resp = &SObjectResponse{}
	err = post(uri, nil, in.(interface{}), resp)

	return
}

func UpdateSObject(id string, in SObject) (err error) {
	uri := strings.Replace(ApiSObjects[in.ApiName()].Urls[rowTemplateKey], idKey, id, 1)

	err = patch(uri, nil, in.(interface{}), nil)

	return
}

func DeleteSObject(id string, in SObject) (err error) {
	uri := strings.Replace(ApiSObjects[in.ApiName()].Urls[rowTemplateKey], idKey, id, 1)

	err = delete(uri, nil)

	return
}

// TODO: Add fields parameter to only retireve needed fields.
func GetSObjectByExternalId(id string, out SObject) (err error) {
	uri := fmt.Sprintf("%v/%v/%v", ApiSObjects[out.ApiName()].Urls[sObjectKey], out.ExternalIdApiName(), id)

	err = get(uri, nil, out.(interface{}))

	return
}

func UpsertSObjectByExternalId(id string, in SObject) (resp *SObjectResponse, err error) {
	uri := fmt.Sprintf("%v/%v/%v", ApiSObjects[in.ApiName()].Urls[sObjectKey], in.ExternalIdApiName(), id)

	resp = &SObjectResponse{}
	err = patch(uri, nil, in.(interface{}), resp)

	return
}

func DeleteSObjectByExternalId(id string, in SObject) (err error) {
	uri := fmt.Sprintf("%v/%v/%v", ApiSObjects[in.ApiName()].Urls[sObjectKey], in.ExternalIdApiName(), id)

	err = delete(uri, nil)

	return
}
