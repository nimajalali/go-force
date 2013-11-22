package sobjects

// Base struct that contains fields that all objects, standard and custom, include.
type BaseSObject struct {
	Attributes       SObjectAttributes `json:"-" force:"attributes,omitempty"`
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

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	// Fractional seconds are handled implicitly by Parse.
	*t, err = Parse(`"2006-01-02T15:04:05.000-0700"`, string(data))
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}
	return []byte(t.Format(`"2006-01-02T15:04:05.000-0700"`)), nil
}
