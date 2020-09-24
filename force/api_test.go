package force

import (
	"bytes"
	"io"
	"testing"
)

func TestForceApi_SetAPIResources(t *testing.T) {

	reader := bytes.NewBufferString(`{
		"tooling": "/services/data/v43.0/tooling",
		"metadata": "/services/data/v43.0/metadata",
		"folders": "/services/data/v43.0/folders",
		"eclair": "/services/data/v43.0/eclair",
		"prechatForms": "/services/data/v43.0/prechatForms",
		"chatter": "/services/data/v43.0/chatter"
	}`)

	type fields struct {
		apiResources map[string]string
	}
	type args struct {
		src io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "setResources",
			fields: fields{
				apiResources: map[string]string{
					"tooling":      "/services/data/v43.0/tooling",
					"metadata":     "/services/data/v43.0/metadata",
					"folders":      "/services/data/v43.0/folders",
					"eclair":       "/services/data/v43.0/eclair",
					"prechatForms": "/services/data/v43.0/prechatForms",
					"chatter":      "/services/data/v43.0/chatter",
				},
			},
			args: args{
				src: reader,
			},
			wantErr: false,
		},
		{
			name: "setResources nil reader",
			fields: fields{
				apiResources: map[string]string{},
			},
			args: args{
				src: nil,
			},
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forceApi := &ForceApi{
				apiResources: make(map[string]string),
				apiSObjects:  make(map[string]*SObjectMetaData),
			}
			if err := forceApi.SetAPIResources(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("ForceApi.SetAPIResources() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(tt.fields.apiResources) != len(forceApi.apiResources) {
				t.Errorf("forceAPI.apiResources length not as expected, wanted: %v, but got: %v", len(forceApi.apiResources), len(tt.fields.apiResources))
			}
			for k, v := range tt.fields.apiResources {
				actual, exists := forceApi.apiResources[k]
				if !exists {
					t.Errorf("Key:%s does not exist in forceApi.apiResources", k)
				}
				if actual != v {
					t.Errorf("Expected value:%s to be associated with key:%s in apiResources. but got: %s instead", v, k, actual)
				}
			}
		})
	}
}

func TestForceApi_SetAPISObjects(t *testing.T) {
	reader := bytes.NewBufferString(validSObjects)
	type fields struct {
		apiSObjects     map[string]*SObjectMetaData
		apiMaxBatchSize int64
	}
	type args struct {
		src io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "setApiObjects",
			fields: fields{
				apiSObjects: map[string]*SObjectMetaData{
					"Applications__c": &SObjectMetaData{},
					"Contact":         &SObjectMetaData{},
					"Liability__c":    &SObjectMetaData{},
					"Other_Income__c": &SObjectMetaData{},
				},
				apiMaxBatchSize: 200,
			},
			args: args{
				src: reader,
			},
			wantErr: false,
		},
		{
			name: "setApiObjects",
			fields: fields{
				apiSObjects:     map[string]*SObjectMetaData{},
				apiMaxBatchSize: 0,
			},
			args: args{
				src: bytes.NewBufferString("{\"malformed\"}"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forceApi := &ForceApi{
				apiResources: make(map[string]string),
				apiSObjects:  make(map[string]*SObjectMetaData),
			}
			if err := forceApi.SetAPISObjects(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("ForceApi.SetAPISObjects() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.fields.apiMaxBatchSize != forceApi.apiMaxBatchSize {
				t.Errorf("Expected ForceApi.apiMaxbatchSize to equal: %v, but got: %v", tt.fields.apiMaxBatchSize, forceApi.apiMaxBatchSize)
			}
			if len(tt.fields.apiSObjects) != len(forceApi.apiSObjects) {
				t.Errorf("Expected ForceApi.apiSObjects to have length: %v, but got: %v", len(tt.fields.apiSObjects), len(forceApi.apiSObjects))
			}
			for k := range tt.fields.apiSObjects {
				_, exists := forceApi.apiSObjects[k]
				if !exists {
					t.Errorf("Expected forceApi.Sobjects to have key:%v", k)
				}
			}
		})
	}
}

const validSObjects = `{
    "encoding": "UTF-8",
    "maxBatchSize": 200,
    "sobjects": [
        {
            "activateable": false,
            "createable": true,
            "custom": true,
            "customSetting": false,
            "deletable": true,
            "deprecatedAndHidden": false,
            "feedEnabled": false,
            "hasSubtypes": false,
            "isSubtype": false,
            "keyPrefix": "a0u",
            "label": "CM Application",
            "labelPlural": "CM Applications",
            "layoutable": true,
            "mergeable": false,
            "mruEnabled": true,
            "name": "Applications__c",
            "queryable": true,
            "replicateable": true,
            "retrieveable": true,
            "searchable": true,
            "triggerable": true,
            "undeletable": true,
            "updateable": true,
            "urls": {
                "compactLayouts": "/services/data/v43.0/sobjects/Applications__c/describe/compactLayouts",
                "rowTemplate": "/services/data/v43.0/sobjects/Applications__c/{ID}",
                "approvalLayouts": "/services/data/v43.0/sobjects/Applications__c/describe/approvalLayouts",
                "defaultValues": "/services/data/v43.0/sobjects/Applications__c/defaultValues?recordTypeId&fields",
                "describe": "/services/data/v43.0/sobjects/Applications__c/describe",
                "quickActions": "/services/data/v43.0/sobjects/Applications__c/quickActions",
                "layouts": "/services/data/v43.0/sobjects/Applications__c/describe/layouts",
                "sobject": "/services/data/v43.0/sobjects/Applications__c"
            }
        },
        {
            "activateable": false,
            "createable": true,
            "custom": false,
            "customSetting": false,
            "deletable": true,
            "deprecatedAndHidden": false,
            "feedEnabled": true,
            "hasSubtypes": false,
            "isSubtype": false,
            "keyPrefix": "003",
            "label": "Contact",
            "labelPlural": "Contacts",
            "layoutable": true,
            "mergeable": true,
            "mruEnabled": true,
            "name": "Contact",
            "queryable": true,
            "replicateable": true,
            "retrieveable": true,
            "searchable": true,
            "triggerable": true,
            "undeletable": true,
            "updateable": true,
            "urls": {
                "compactLayouts": "/services/data/v43.0/sobjects/Contact/describe/compactLayouts",
                "rowTemplate": "/services/data/v43.0/sobjects/Contact/{ID}",
                "approvalLayouts": "/services/data/v43.0/sobjects/Contact/describe/approvalLayouts",
                "defaultValues": "/services/data/v43.0/sobjects/Contact/defaultValues?recordTypeId&fields",
                "listviews": "/services/data/v43.0/sobjects/Contact/listviews",
                "describe": "/services/data/v43.0/sobjects/Contact/describe",
                "quickActions": "/services/data/v43.0/sobjects/Contact/quickActions",
                "layouts": "/services/data/v43.0/sobjects/Contact/describe/layouts",
                "sobject": "/services/data/v43.0/sobjects/Contact"
            }
        },
        {
            "activateable": false,
            "createable": true,
            "custom": true,
            "customSetting": false,
            "deletable": true,
            "deprecatedAndHidden": false,
            "feedEnabled": false,
            "hasSubtypes": false,
            "isSubtype": false,
            "keyPrefix": "aCl",
            "label": "Liability",
            "labelPlural": "Liabilities",
            "layoutable": true,
            "mergeable": false,
            "mruEnabled": false,
            "name": "Liability__c",
            "queryable": true,
            "replicateable": true,
            "retrieveable": true,
            "searchable": true,
            "triggerable": true,
            "undeletable": true,
            "updateable": true,
            "urls": {
                "compactLayouts": "/services/data/v43.0/sobjects/Liability__c/describe/compactLayouts",
                "rowTemplate": "/services/data/v43.0/sobjects/Liability__c/{ID}",
                "approvalLayouts": "/services/data/v43.0/sobjects/Liability__c/describe/approvalLayouts",
                "defaultValues": "/services/data/v43.0/sobjects/Liability__c/defaultValues?recordTypeId&fields",
                "describe": "/services/data/v43.0/sobjects/Liability__c/describe",
                "quickActions": "/services/data/v43.0/sobjects/Liability__c/quickActions",
                "layouts": "/services/data/v43.0/sobjects/Liability__c/describe/layouts",
                "sobject": "/services/data/v43.0/sobjects/Liability__c"
            }
        },
        {
            "activateable": false,
            "createable": true,
            "custom": true,
            "customSetting": false,
            "deletable": true,
            "deprecatedAndHidden": false,
            "feedEnabled": false,
            "hasSubtypes": false,
            "isSubtype": false,
            "keyPrefix": "aCm",
            "label": "Other Income",
            "labelPlural": "Other Income",
            "layoutable": true,
            "mergeable": false,
            "mruEnabled": false,
            "name": "Other_Income__c",
            "queryable": true,
            "replicateable": true,
            "retrieveable": true,
            "searchable": true,
            "triggerable": true,
            "undeletable": true,
            "updateable": true,
            "urls": {
                "compactLayouts": "/services/data/v43.0/sobjects/Other_Income__c/describe/compactLayouts",
                "rowTemplate": "/services/data/v43.0/sobjects/Other_Income__c/{ID}",
                "approvalLayouts": "/services/data/v43.0/sobjects/Other_Income__c/describe/approvalLayouts",
                "defaultValues": "/services/data/v43.0/sobjects/Other_Income__c/defaultValues?recordTypeId&fields",
                "describe": "/services/data/v43.0/sobjects/Other_Income__c/describe",
                "quickActions": "/services/data/v43.0/sobjects/Other_Income__c/quickActions",
                "layouts": "/services/data/v43.0/sobjects/Other_Income__c/describe/layouts",
                "sobject": "/services/data/v43.0/sobjects/Other_Income__c"
            }
        }
    ]
}`
