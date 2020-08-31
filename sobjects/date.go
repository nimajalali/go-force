package sobjects

import (
	"fmt"
	"github.com/nimajalali/go-force/forcejson"
	"time"
)

const SFTIMEFORMAT1 = "2006-01-02T15:04:05.000-0700"
const SFTIMEFORMAT2 = "2006-01-02T15:04:05.999Z"
const SFTIMEFORMAT3 = "2006-01-02"

var sfdcMinDate = time.Date(1700, 1, 1, 0, 0, 0, 0, time.UTC)

// Represents a SFDC Date. Implements marshaling / unmarshaling as a Go Time.
type Time time.Time

// Convenience Go time.Time constructor.
func AsTime(t time.Time) *Time {
	if t.IsZero() || t.Before(sfdcMinDate) {
		return nil
	}
	ret := Time(t.UTC().Truncate(time.Second)) // SFDC doesn't store ms
	return &ret
}

func ParseTime(str string) (*Time, error) {
	tm, err := time.Parse(SFTIMEFORMAT1, str)
	if err != nil {
		tm, err = time.Parse(SFTIMEFORMAT2, str)
	}
	if err != nil {
		tm, err = time.Parse(SFTIMEFORMAT3, str)
	}
	if err != nil {
		return nil, err
	}
	return AsTime(tm), nil
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in Salesforce time format.
func (t Time) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte{'"', '"'}, nil
	}
	b := make([]byte, 0, len(SFTIMEFORMAT1)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, SFTIMEFORMAT1)
	b = append(b, '"')
	return b, nil
}

// Convenience Go time.Time converstion.
func (t *Time) Time() time.Time {
	if t == nil {
		return time.Time{}
	}
	return time.Time(*t)
}

// Convenience Stringer implementation.
func (t Time) String() string {
	if time.Time(t).IsZero() {
		return ""
	}
	return time.Time(t).Format(SFTIMEFORMAT1)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in Salesforce time format.
func (t *Time) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid JSON string")
	}

	str := string(data[1 : len(data)-1])
	if len(str) == 0 {
		*t = Time(time.Time{})
		return nil
	}

	tm, err := ParseTime(str)
	if err != nil {
		return err
	}
	*t = *tm
	return nil
}

var _ forcejson.Unmarshaler = (*Time)(nil)
var _ forcejson.Marshaler = Time{}
