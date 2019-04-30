package sobjects

import (
	"reflect"
	"strings"
)

var baseFieldNameMap map[string]string

func init() {
	baseFieldNameMap = map[string]string{
		"Id":               "Id",
		"IsDeleted":        "IsDeleted",
		"Name":             "Name",
		"CreatedDate":      "CreatedDate",
		"CreatedById":      "CreatedById",
		"LastModifiedDate": "LastModifiedDate",
		"LastModifiedById": "LastModifiedById",
		"SystemModstamp":   "SystemModstamp",
	}
}

// Base struct that contains fields that all objects, standard and custom, include.
type BaseSObject struct {
	Attributes       SObjectAttributes `force:"attributes,omitempty" json:"-"`
	Id               string            `force:",omitempty" json:",omitempty"`
	IsDeleted        bool              `force:",omitempty" json:",omitempty"`
	Name             string            `force:",omitempty" json:",omitempty"`
	CreatedDate      *Time             `force:",omitempty" json:",omitempty"`
	CreatedById      string            `force:",omitempty" json:",omitempty"`
	LastModifiedDate *Time             `force:",omitempty" json:",omitempty"`
	LastModifiedById string            `force:",omitempty" json:",omitempty"`
	SystemModstamp   string            `force:",omitempty" json:",omitempty"`
}

type SObjectAttributes struct {
	Type string `force:"type,omitempty"`
	Url  string `force:"url,omitempty"`
}

// Implementing this here because most objects don't have an external id and as such this is not needed.
// Feel free to override this function when embedding the BaseSObject in other structs.
func (b BaseSObject) ExternalIdApiName() string {
	return ""
}

// Fields that are returned in every query response. Use this to build custom structs.
// type MyCustomQueryResponse struct {
// 	BaseQuery
// 	Records []sobjects.Account `json:"records" force:"records"`
// }
type BaseQuery struct {
	Done           bool    `json:"Done" force:"done"`
	TotalSize      float64 `json:"TotalSize" force:"totalSize"`
	NextRecordsUri string  `json:"NextRecordsUrl" force:"nextRecordsUrl"`
}

// ConvertFieldNames takes in any interface that inplements SObject and a comma separated list of json field names.
// It converts the json field names to the force struct tag stated equivalent.
func ConvertFieldNames(obj interface{}, jsonFields string) string {
	if jsonFields != "" {
		fields := strings.Split(jsonFields, ",")

		length := len(fields)
		if length > 0 {
			mapping := fieldNameMapping(obj)

			var forceFields []string
			for _, field := range fields {
				if forceField, ok := mapping[field]; ok {
					forceFields = append(forceFields, forceField)
				}
			}

			return strings.Join(forceFields, ",")
		}
	}

	return ""
}

// Helper function used in ConvertFieldNames
func fieldNameMapping(obj interface{}) map[string]string {
	st := reflect.TypeOf(obj)
	fl := st.NumField()

	jsonToForce := make(map[string]string, fl)

	for i := 0; i < fl; i++ {
		sf := st.Field(i)
		jName := strings.SplitN(sf.Tag.Get("json"), ",", 2)[0]
		fName := strings.SplitN(sf.Tag.Get("force"), ",", 2)[0]

		if jName == "-" {
			continue
		}

		if fName == "-" {
			continue
		}

		if jName == "" {
			jName = sf.Name
		}

		if fName == "" {
			fName = sf.Name
		}

		jsonToForce[jName] = fName
	}

	for k, v := range baseFieldNameMap {
		jsonToForce[k] = v
	}

	return jsonToForce
}
