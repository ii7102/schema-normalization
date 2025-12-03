package schema

import (
	"time"

	"github.com/ii7102/schema-normalization/errors"
)

func matchesBaseType(baseType BaseType, value any) bool {
	switch value.(type) {
	case bool:
		return baseType == Boolean
	case int, int8, int16, int32, int64:
		return baseType == Integer
	case string:
		return baseType == String
	case float32, float64:
		return baseType == Float
	default:
		return false
	}
}

func validateEnumValues(baseType BaseType, enumValues ...any) error {
	switch baseType {
	case Date, Timestamp, DateTime, Object:
		return errors.WrappedError(errEnumsOfBaseTypeAreNotSupported, "enums of %s are not supported", baseType)
	default:
	}

	if len(enumValues) == 0 {
		return errEnumValuesCannotBeNilOrEmpty
	}

	for _, enumValue := range enumValues {
		if enumValue == nil {
			return errEnumValueCannotBeNil
		}

		if !matchesBaseType(baseType, enumValue) {
			return errors.WrappedError(errEnumValueDoesNotMatchBaseType, "enum value %v is not of %s type", enumValue, baseType)
		}
	}

	return nil
}

func validateObjectFields(objectFields map[Field]FieldType) error {
	if objectFields == nil {
		return errObjectFieldsCannotBeNil
	}

	if len(objectFields) == 0 {
		return errObjectFieldsCannotBeEmpty
	}

	for field, fieldType := range objectFields {
		if fieldType.baseType == Object {
			return errors.WrappedError(errNestedObjectsAreNotSupported, "object field '%s' cannot be of type object ", field)
		}
	}

	return nil
}

func validateDateTime(layout string, baseType BaseType) error {
	validateLayoutFunc := validateDateTimeFuncs(baseType)

	if validateLayoutFunc == nil {
		return errors.WrappedError(errBaseTypeWithoutLayout, "base type %s does not have a layout", baseType)
	}

	return validateLayoutFunc(layout)
}

func validateDateTimeFuncs(baseType BaseType) func(string) error {
	switch baseType {
	case Timestamp:
		return validateTimestampLayout
	case Date:
		return validateDateLayout
	case DateTime:
		return validateDateTimeLayout
	default:
		return nil
	}
}

func validateDateLayout(_ string) error {
	return nil
}

func validateTimestampLayout(layout string) error {
	switch layout {
	case time.TimeOnly, time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano:
		return nil
	default:
		return errInvalidTimestampLayout
	}
}

func validateDateTimeLayout(layout string) error {
	switch layout {
	case time.RFC3339, time.DateTime, time.RFC3339Nano, time.RFC1123, time.RFC1123Z,
		time.RFC850, time.RFC822, time.RFC822Z, time.ANSIC, time.UnixDate, time.RubyDate:
		return nil
	default:
		return errInvalidDateTimeLayout
	}
}
