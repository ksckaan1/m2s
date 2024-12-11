package mpftostruct

import "errors"

var (
	ErrValueMustBePointer = errors.New("value must be a pointer")
	ErrValueCannotBeNil   = errors.New("value cannot be nil")
	ErrValueMustBeStruct  = errors.New("value must be a struct")
	ErrInvalidFieldType   = errors.New("invalid field type")
)
