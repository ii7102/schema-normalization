package rules

import (
	"errors"
	"fmt"
)

func validateEnumValues(baseType BaseType, enumValues ...any) error {
	switch baseType {
	case Date, Timestamp, DateTime, Object:
		return fmt.Errorf("enums of %s are not supported", baseType)
	}

	if len(enumValues) == 0 {
		return errors.New("enum values cannot be nil or empty")
	}

	for _, enumValue := range enumValues {
		if enumValue == nil {
			return errors.New("enum value cannot be nil")
		}

		if !matchesBaseType(baseType, enumValue) {
			return fmt.Errorf("enum value %v is not of %s type", enumValue, baseType)
		}
	}

	return nil
}

func matchesBaseType(baseType BaseType, value any) bool {
	switch baseType {
	case Boolean:
		_, ok := value.(bool)

		return ok
	case Integer:
		switch value.(type) {
		case int, int8, int16, int32, int64:
			return true
		default:
			return false
		}
	case String:
		_, ok := value.(string)

		return ok
	case Float:
		switch value.(type) {
		case float32, float64:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func validateObjectFields(objectFields map[Field]FieldType) error {
	if objectFields == nil {
		return errors.New("object fields cannot be nil")
	}

	if len(objectFields) == 0 {
		return errors.New("object fields cannot be empty")
	}

	for field, fieldType := range objectFields {
		if fieldType.baseType == Object {
			return fmt.Errorf("object field '%s' cannot be of type object (nested objects are not supported)", field)
		}
	}

	return nil
}
