package rules

import (
	"errors"
	"fmt"
)

var (
	errEnumValuesNotSet               = errors.New("enum values are not set")
	errEnumValueCannotBeNil           = errors.New("enum value cannot be nil")
	errEnumValuesCannotBeNilOrEmpty   = errors.New("enum values cannot be nil or empty")
	errEnumValueDoesNotMatchBaseType  = errors.New("enum value does not match base type")
	errEnumsOfBaseTypeAreNotSupported = errors.New("enums of base type are not supported")

	errObjectFieldsCannotBeNil      = errors.New("object fields cannot be nil")
	errObjectFieldsCannotBeEmpty    = errors.New("object fields cannot be empty")
	errNestedObjectsAreNotSupported = errors.New("nested objects are not supported")

	errExpectedOutputMismatch = errors.New("expected output does not match actual output")
)

// WrappedError wraps the given error with the given message and arguments.
func WrappedError(err error, msg string, args ...any) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %v", err, fmt.Sprintf(msg, args...))
}
