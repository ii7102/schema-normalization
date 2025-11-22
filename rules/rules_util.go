package rules

func validateEnumValues(baseType BaseType, enumValues ...any) error {
	switch baseType {
	case Date, Timestamp, DateTime, Object:
		return WrappedError(errEnumsOfBaseTypeAreNotSupported, "enums of %s are not supported", baseType)
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
			return WrappedError(errEnumValueDoesNotMatchBaseType, "enum value %v is not of %s type", enumValue, baseType)
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
		return errObjectFieldsCannotBeNil
	}

	if len(objectFields) == 0 {
		return errObjectFieldsCannotBeEmpty
	}

	for field, fieldType := range objectFields {
		if fieldType.baseType == Object {
			return WrappedError(errNestedObjectsAreNotSupported, "object field '%s' cannot be of type object ", field)
		}
	}

	return nil
}
