package gonormalizer

import "errors"

var (
	errNormalizedValueIsNotAnArray = errors.New("normalized value is not an array")

	errInvalidValue        = errors.New("invalid value")
	errInvalidEnumValue    = errors.New("invalid enum value")
	errInvalidObjectValue  = errors.New("invalid object value")
	errInvalidArrayValue   = errors.New("invalid array value")
	errInvalidBooleanValue = errors.New("invalid boolean value")
	errInvalidIntegerValue = errors.New("invalid integer value")
	errInvalidFloatValue   = errors.New("invalid float value")
)
