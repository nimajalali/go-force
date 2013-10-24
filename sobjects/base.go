package sobjects

// Base struct that contains fields that all objects, standard and custom, include.
type BaseSObject struct {
	Attributes       SObjectAttributes `force:"attributes,omitempty"`
	Id               string            `force:",omitempty"`
	IsDeleted        bool              `force:",omitempty"`
	Name             string            `force:",omitempty"`
	CreatedDate      string            `force:",omitempty"`
	CreatedById      string            `force:",omitempty"`
	LastModifiedDate string            `force:",omitempty"`
	LastModifiedById string            `force:",omitempty"`
	SystemModstamp   string            `force:",omitempty"`
}

type SObjectAttributes struct {
	Type string `force:"type,omitempty"`
	Url  string `force:"url,omitempty"`
}

// Implementing this here because most object don't have an external id and as such this is not needed.
func (b *BaseSObject) ExternalIdApiName() string {
	return ""
}
