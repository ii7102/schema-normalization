package schema

import (
	"errors"
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
	errObjectFieldsNotSet           = errors.New("object fields are not set")
	errFieldTypeIsNotAnObject       = errors.New("field type is not an object")

	errInvalidTimestampLayout = errors.New("invalid timestamp layout")
	errInvalidDateTimeLayout  = errors.New("invalid dateTime layout")
	errBaseTypeWithoutLayout  = errors.New("base type without layout")
)
