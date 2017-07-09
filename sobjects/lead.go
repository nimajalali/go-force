package sobjects

type Lead struct {
	BaseSObject
	Company       string `force:",omitempty"`
	ConvertedDate string `force:",omitempty"`
	FirstName     string `force:",omitempty"`
	IsConverted   bool   `force:",omitempty"`
	IsDeleted     bool   `force:",omitempty"`
	LastName      string `force:",omitempty"`
	OwnerId       string `force:",omitempty"`
	Status        string `force:",omitempty"`
}

func (t *Lead) APIName() string {
	return "Lead"
}

type LeadQueryResponse struct {
	BaseQuery
	Records []Lead `json:"Records" force:"records"`
}
