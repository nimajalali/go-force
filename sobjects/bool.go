package sobjects

import "encoding/json"

// Don't use this! It was an interesting effort but in reality all you need is a ptr to a bool. *bool will solve all your problems. :)
// Used to represent empty bools. Go types are always instantiated with a default value, for bool the default value is false.
// This makes it difficult to update an SObject without overwriting any boolean field to false.
// This package solves the issue by representing a bool as an int and implementing the marshal/unmarshal json interface.
// It is a drop in replacement for the bool type.
// 1 is true
// 0 is nil
// -1 is false
// Unmarshalling: false will be unmarshalled to -1, true will be unmarshalled to 1
// If no value is set the unmarshaller will skip the field and the int will default to 0.
// Marshalling: -1 will be marshaled to false, 1 will be marshaled to true, and
// 0 will be marshaled to nothing (assuming the field has the omitempty json tag `json:",omitempty"`)
type SFBool int

func (t *SFBool) MarshalJSON() ([]byte, error) {
	if *t == 1 {
		return json.Marshal(true)
	} else if *t == -1 {
		return json.Marshal(false)
	}
	return json.Marshal(0)
}

func (t *SFBool) UnmarshalJSON(data []byte) error {
	b := string(data)
	if b == "true" {
		*t = 1
	} else if b == "false" {
		*t = -1
	}
	return nil
}

func (t *SFBool) Bool() bool {
	if *t == 1 {
		return true
	}
	return false
}
