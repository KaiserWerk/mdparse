package mdparse

import "reflect"

// An InvalidUnmarshalError describes an invalid argument.
// The argument must be a non-nil pointer.
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "mdparse: Read(nil)"
	}
	if e.Type.Kind() != reflect.Pointer {
		return "mdparse: Read(non-pointer " + e.Type.String() + ")"
	}
	return "mdparse: Read(nil " + e.Type.String() + ")"
}
