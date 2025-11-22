package gonormalizer

import (
	"fmt"
	"math"
	"reflect"
	"slices"
	"strconv"
	"time"

	"diploma/rules"
)

var _ rules.AbstractNormalizer = (*normalizer)(nil)

type normalizer struct {
	*rules.BaseNormalizer
}

// NewNormalizer creates a new Go normalizer with the given options.
func NewNormalizer(options ...rules.NormalizerOption) (*normalizer, error) {
	base, err := rules.NewBaseNormalizer(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create base normalizer: %w", err)
	}

	return &normalizer{BaseNormalizer: base}, nil
}

func (n normalizer) Normalize(data map[string]any) (map[string]any, error) {
	normalizedData := make(map[string]any, len(data))
	for key, value := range data {
		fieldType, ok := n.Fields[rules.Field(key)]
		if !ok {
			continue
		}

		normalizedValue, err := normalize(value, fieldType)
		if err != nil {
			return nil, err
		}

		normalizedData[key] = normalizedValue
	}

	return normalizedData, nil
}

func normalize(value any, fieldType rules.FieldType) (any, error) {
	var (
		normalizedValue any
		err             error
	)

	if value == nil {
		return value, nil
	}

	if fieldType.IsArray() {
		normalizedValue, err = normalizeArray(value, fieldType)
	} else {
		normalizedValue, err = normalizeValue(value, fieldType)
	}

	if err != nil {
		return nil, err
	}

	if len(fieldType.EnumValues()) == 0 {
		return normalizedValue, nil
	}

	if normalizedValue == nil {
		return nil, rules.WrappedError(errInvalidEnumValue, "value is nil")
	}

	if err = validateEnums(normalizedValue, fieldType.EnumValues(), fieldType.IsArray()); err != nil {
		return nil, err
	}

	return normalizedValue, nil
}

func normalizeArray(value any, fieldType rules.FieldType) (any, error) {
	reflectVal := reflect.ValueOf(value)
	if reflectVal.Kind() != reflect.Slice {
		return nil, rules.WrappedError(errInvalidArrayValue, "value: %v", value)
	}

	normalizedArray := make([]any, 0, reflectVal.Len())
	for i := range reflectVal.Len() {
		elem := reflectVal.Index(i).Interface()
		if elem == nil {
			continue
		}

		normalizedValue, err := normalizeValue(elem, fieldType)
		if err != nil {
			return nil, err
		}

		normalizedArray = append(normalizedArray, normalizedValue)
	}

	return normalizedArray, nil
}

func normalizeValue(value any, fieldType rules.FieldType) (any, error) {
	switch fieldType.BaseType() {
	case rules.Boolean:
		return normalizeBoolean(value)
	case rules.Integer:
		return normalizeInteger(value)
	case rules.String:
		return normalizeString(value)
	case rules.Float:
		return normalizeFloat(value)
	case rules.Date, rules.Timestamp, rules.DateTime:
		return normalizeDateTime(value, fieldType.BaseType())
	case rules.Object:
		return normalizeObject(value, fieldType)
	default:
		return nil, rules.WrappedError(errInvalidValue, "%v, type: %T", value, value)
	}
}

func normalizeBoolean(value any) (any, error) {
	switch value := value.(type) {
	case bool:
		return value, nil
	case string:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue, nil
		}

		return nil, rules.WrappedError(errInvalidBooleanValue, "string value: %s cannot be parsed to boolean", value)
	}

	return nil, rules.WrappedError(errInvalidBooleanValue, "value: %v", value)
}

func normalizeInteger(value any) (any, error) {
	switch value := value.(type) {
	case int, int8, int16, int32, int64:
		return value, nil
	case float32:
		return int64(value), nil
	case float64:
		return int64(value), nil
	case string:
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue, nil
		}

		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return int(math.Round(floatValue)), nil
		}

		return nil, rules.WrappedError(errInvalidIntegerValue, "string value: %s cannot be parsed to number", value)
	}

	return nil, rules.WrappedError(errInvalidIntegerValue, "value: %v, type: %T", value, value)
}

func normalizeString(value any) (any, error) {
	return fmt.Sprintf("%v", value), nil
}

func normalizeFloat(value any) (any, error) {
	switch value := value.(type) {
	case int:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case float32:
		return float64(value), nil
	case float64:
		return value, nil
	case string:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue, nil
		}

		return nil, rules.WrappedError(errInvalidFloatValue, "string value: %s cannot be parsed to float", value)
	default:
		return nil, rules.WrappedError(errInvalidFloatValue, "value: %v, type: %T", value, value)
	}
}

func normalizeDateTime(value any, baseType rules.BaseType) (any, error) {
	stringValue, ok := value.(string)
	if !ok {
		return nil, rules.WrappedError(errInvalidValue, "%s value: %v", baseType, value)
	}

	layoutFormats := layoutFormats(baseType)
	for _, layoutFormat := range layoutFormats {
		parsedTime, err := time.Parse(layoutFormat, stringValue)
		if err == nil {
			return parsedTime.Format(baseLayoutFormat(baseType)), nil
		}
	}

	return nil, rules.WrappedError(errInvalidValue, "%s value: string %s cannot be formatted", baseType, stringValue)
}

func normalizeObject(value any, fieldType rules.FieldType) (any, error) {
	valueMap, ok := value.(map[string]any)
	if !ok {
		return nil, rules.WrappedError(errInvalidObjectValue, "value: %v", value)
	}

	normalizedObject := make(map[string]any, len(valueMap))
	for key, val := range valueMap {
		objectField, ok := fieldType.ObjectFields()[rules.Field(key)]
		if !ok {
			continue
		}

		normalizedValue, err := normalize(val, objectField)
		if err != nil {
			return nil, err
		}

		normalizedObject[key] = normalizedValue
	}

	return normalizedObject, nil
}

func validateEnums(normalizedValue any, enumValues []any, isArray bool) error {
	if !isArray {
		if !slices.Contains(enumValues, normalizedValue) {
			return rules.WrappedError(errInvalidEnumValue, "value: %v", normalizedValue)
		}

		return nil
	}

	normalizedValueArray, ok := normalizedValue.([]any)
	// This shouldn't happen as previous normalization checks for this, but none the less we should check it
	if !ok {
		return errNormalizedValueIsNotAnArray
	}

	for _, value := range normalizedValueArray {
		err := validateEnums(value, enumValues, false)
		if err != nil {
			return err
		}
	}

	return nil
}
