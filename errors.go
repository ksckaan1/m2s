package m2s

import "errors"

var (
	ErrValueMustBePointer = errors.New("value must be a pointer")
	ErrValueCannotBeNil   = errors.New("value cannot be nil")
	ErrValueMustBeStruct  = errors.New("value must be a struct")
	ErrInvalidFieldType   = errors.New("invalid field type")
)

type ErrParseFailed struct {
	Field string
	Err   error
}

func (e ErrParseFailed) Error() string {
	return "failed to parse field " + e.Field + ": " + e.Err.Error()
}
