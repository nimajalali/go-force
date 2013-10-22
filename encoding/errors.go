package encoding

import (
	"reflect"
	"strconv"
)

// An UnsupportedTypeError is returned by Marshal when attempting
// to encode an unsupported value type.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "force: unsupported type: " + e.Type.String()
}

type UnsupportedValueError struct {
	Value reflect.Value
	Str   string
}

func (e *UnsupportedValueError) Error() string {
	return "force: unsupported value: " + e.Str
}

// An DecodeTypeError describes a value that was
// not appropriate for a value of a specific Go type.
type DecodeTypeError struct {
	Value string       // description of value - "bool", "array", "number -5"
	Type  reflect.Type // type of Go value it could not be assigned to
}

func (e *DecodeTypeError) Error() string {
	return "force: cannot decode " + e.Value + " into Go value of type " + e.Type.String()
}

// An DecodeFieldError describes a object key that
// led to an unexported (and therefore unwritable) struct field.
// (No longer used; kept for compatibility.)
type DecodeFieldError struct {
	Key   string
	Type  reflect.Type
	Field reflect.StructField
}

func (e *DecodeFieldError) Error() string {
	return "force: cannot decode object key " + strconv.Quote(e.Key) + " into unexported field " + e.Field.Name + " of type " + e.Type.String()
}

// An InvalidDecodeError describes an invalid argument passed to Decode.
// (The argument to Decode must be a non-nil pointer.)
type InvalidDecodeError struct {
	Type reflect.Type
}

func (e *InvalidDecodeError) Error() string {
	if e.Type == nil {
		return "force: Decode(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "force: Decode(non-pointer " + e.Type.String() + ")"
	}
	return "force: Decode(nil " + e.Type.String() + ")"
}
