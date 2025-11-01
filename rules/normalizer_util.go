package rules

func validateEnumValues(baseType BaseType, enumValues []any) (invalidEnumValue any) {
	for _, v := range enumValues {
		if v == nil {
			continue
		}
		if !matchesBaseType(baseType, v) {
			return v
		}
	}
	return nil
}

func matchesBaseType(baseType BaseType, v any) bool {
	switch baseType {
	case Boolean:
		_, ok := v.(bool)
		return ok
	case Integer:
		switch v.(type) {
		case int, int8, int16, int32, int64:
			return true
		}
	case String:
		_, ok := v.(string)
		return ok
	case Float:
		switch v.(type) {
		case float32, float64:
			return true
		}
	}
	return false
}
