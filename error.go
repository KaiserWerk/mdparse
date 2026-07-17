package mdparse

import "reflect"

// An InvalidUnmarshalError describes an invalid argument passed to [Parse].
// (The argument to [Parse] must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "mdparse: Parse(nil)"
	}
	if e.Type.Kind() != reflect.Pointer {
		return "mdparse: Parse(non-pointer " + e.Type.String() + ")"
	}
	return "mdparse: Parse(nil " + e.Type.String() + ")"
}
