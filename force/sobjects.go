package force

import (
	"fmt"
)

const (
	sobjectsResourceKey = "sobjects"
)

// Interface needed to expose the api name of the sobject. Needed to generate the uri.
type SObject interface {
	ApiName() string
}

// Response recieved from force.com API after insert, update, or delete of an sobject.
type SObjectResponse struct {
	Id      string    `force:"id,omitempty"`
	Errors  ApiErrors `force:"error,omitempty"`
	Success bool      `force:"success,omitempty"`
}

// TODO: Add fields parameter to only retireve needed fields.
func GetSObject(id string, out SObject) (err error) {
	uri := fmt.Sprintf("%v/%v/%v", ApiResources[sobjectsResourceKey], out.ApiName(), id)

	err = get(uri, nil, out.(interface{}))

	return
}

func InsertSObject(in SObject) (resp *SObjectResponse, err error) {
	uri := fmt.Sprintf("%v/%v", ApiResources[sobjectsResourceKey], in.ApiName())

	resp = &SObjectResponse{}
	err = post(uri, nil, in.(interface{}), resp)

	return
}

func UpdateSObject(id string, in SObject) (err error) {
	uri := fmt.Sprintf("%v/%v/%v", ApiResources[sobjectsResourceKey], in.ApiName(), id)

	err = patch(uri, nil, in.(interface{}))

	return
}
