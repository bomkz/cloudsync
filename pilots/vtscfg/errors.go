package vtscfg

import (
	"fmt"
	"reflect"
)

// SyntaxError reports a parse error together with the 1-based source line
// it occurred on. (encoding/json reports a byte offset instead; this
// format is line-oriented, so a line number is more useful.)
type SyntaxError struct {
	Line int
	Msg  string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("vtscfg: syntax error on line %d: %s", e.Line, e.Msg)
}

// InvalidUnmarshalError mirrors json.InvalidUnmarshalError: returned when
// Unmarshal/Decode is called with something other than a non-nil pointer.
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "vtscfg: Unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Ptr {
		return "vtscfg: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "vtscfg: Unmarshal(nil " + e.Type.String() + ")"
}

// UnmarshalTypeError mirrors json.UnmarshalTypeError: the value in the
// file couldn't be assigned to the requested Go type.
type UnmarshalTypeError struct {
	Value string       // raw text (or "object"/"tuple") that couldn't be used
	Type  reflect.Type // Go type it couldn't be assigned to
	Field string       // full path of the field, if known
}

func (e *UnmarshalTypeError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("vtscfg: cannot unmarshal %q into Go struct field %s of type %s", e.Value, e.Field, e.Type)
	}
	return fmt.Sprintf("vtscfg: cannot unmarshal %q into Go value of type %s", e.Value, e.Type)
}

// UnsupportedTypeError mirrors json.UnsupportedTypeError: Marshal was
// asked to encode a Go type it has no representation for (channels,
// funcs, complex numbers, ...).
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "vtscfg: unsupported type: " + e.Type.String()
}

// Disclaimer: All code in this document has been taken from LLM outputs.
// TODO: Consider rewriting eventually.
