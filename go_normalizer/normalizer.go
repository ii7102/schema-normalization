package gonormalizer

import (
	"fmt"
	"math"
	"reflect"
	"slices"
	"strconv"
	"strings"
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
		return nil, rules.WrappedError(err, "failed to create base normalizer")
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

func normalize(value any, fieldType rules.FieldType) (normalizedValue any, err error) {
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

	if err = validateEnum(normalizedValue, fieldType.EnumValues(), fieldType.IsArray()); err != nil {
		return nil, err
	}

	return normalizedValue, nil
}

func normalizeArray(value any, fieldType rules.FieldType) (any, error) {
	reflectVal := reflect.ValueOf(value)
	if reflectVal.Kind() != reflect.Slice {
		return nil, rules.WrappedError(errInvalidArrayValue, "%v", value)
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
		return normalizeString(value), nil
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

func normalizeBoolean(value any) (bool, error) {
	switch value := value.(type) {
	case bool:
		return value, nil
	case string:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return false, rules.WrappedError(errInvalidBooleanValue, "string value: %s, error: %v", value, err)
		}

		return boolValue, nil
	default:
		return false, rules.WrappedError(errInvalidBooleanValue, "%v, type: %T", value, value)
	}
}

func normalizeNumberFromString(value string) (float64, error) {
	if strings.Contains(value, ".") {
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, fmt.Errorf("string value: %s, error: %w", value, err)
		}

		return floatValue, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("string value: %s, error: %w", value, err)
	}

	return float64(intValue), nil
}

func normalizeInteger(value any) (int64, error) {
	reflectVal := reflect.ValueOf(value)
	switch reflectVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflectVal.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue := reflectVal.Uint()
		if uintValue > math.MaxInt64 {
			return 0, rules.WrappedError(errInvalidIntegerValue, "uint64 value: %d is greater than max int64 value", uintValue)
		}

		return int64(uintValue), nil
	case reflect.Float32, reflect.Float64:
		return int64(reflectVal.Float()), nil
	case reflect.String:
		value, err := normalizeNumberFromString(reflectVal.String())
		if err != nil {
			return 0, rules.WrappedError(errInvalidIntegerValue, "%s", err)
		}

		return int64(value), nil
	default:
		return 0, rules.WrappedError(errInvalidIntegerValue, "%v, type: %T", value, value)
	}
}

func normalizeString(value any) string {
	return fmt.Sprintf("%v", value)
}

func normalizeFloat(value any) (float64, error) {
	reflectVal := reflect.ValueOf(value)
	switch reflectVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(reflectVal.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(reflectVal.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return reflectVal.Float(), nil
	case reflect.String:
		value, err := normalizeNumberFromString(reflectVal.String())
		if err != nil {
			return 0, rules.WrappedError(errInvalidFloatValue, "%s", err)
		}

		return value, nil
	default:
		return 0, rules.WrappedError(errInvalidFloatValue, "%v, type: %T", value, value)
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
		return nil, rules.WrappedError(errInvalidObjectValue, "%v", value)
	}

	objectFields := fieldType.ObjectFields()

	normalizedObject := make(map[string]any, len(objectFields))
	for key, val := range valueMap {
		objectField, ok := objectFields[rules.Field(key)]
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

func validateEnum(normalizedValue any, enumValues []any, isArray bool) error {
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
		err := validateEnum(value, enumValues, false)
		if err != nil {
			return err
		}
	}

	return nil
}
