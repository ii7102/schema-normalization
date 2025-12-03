package gonormalizer

import (
	"fmt"
	"math"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/ii7102/schema-normalization/errors"
	"github.com/ii7102/schema-normalization/schema"
)

var _ schema.AbstractNormalizer = (*normalizer)(nil)

type normalizer struct {
	*schema.BaseNormalizer
}

// NewNormalizer creates a new Go normalizer with the given options.
func NewNormalizer(opts ...schema.NormalizerOption) (*normalizer, error) {
	base, err := schema.NewBaseNormalizer(opts...)
	if err != nil {
		return nil, errors.WrappedError(err, "failed to create base normalizer")
	}

	return &normalizer{BaseNormalizer: base}, nil
}

func (n normalizer) Normalize(data map[string]any) (map[string]any, error) {
	normalizedData := make(map[string]any, len(data))
	for key, value := range data {
		field := schema.Field(key)
		if !n.HasField(field) {
			continue
		}

		normalizedValue, err := normalize(value, n.FieldType(field))
		if err != nil {
			return nil, err
		}

		normalizedData[key] = normalizedValue
	}

	return normalizedData, nil
}

func normalize(value any, fieldType schema.FieldType) (any, error) {
	if value == nil {
		return value, nil
	}

	normalizedValue, err := normalizeArray(value, fieldType)
	if err != nil {
		return nil, err
	}

	if len(fieldType.EnumValues()) == 0 {
		return normalizedValue, nil
	}

	if err = validateEnum(normalizedValue, fieldType.EnumValues(), fieldType.IsArray()); err != nil {
		return nil, err
	}

	return normalizedValue, nil
}

func normalizeArray(value any, fieldType schema.FieldType) (any, error) {
	if !fieldType.IsArray() {
		return normalizeValue(value, fieldType)
	}

	reflectVal := reflect.ValueOf(value)
	if reflectVal.Kind() != reflect.Slice {
		return nil, errors.WrappedError(errInvalidArrayValue, "%v", value)
	}

	normalizedArray := make([]any, 0, reflectVal.Len())
	for i := range reflectVal.Len() {
		elem := reflectVal.Index(i).Interface()
		if elem == nil {
			normalizedArray = append(normalizedArray, nil)

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

func normalizeValue(value any, fieldType schema.FieldType) (any, error) {
	switch fieldType.BaseType() {
	case schema.Boolean:
		return normalizeBoolean(value)
	case schema.Integer:
		return normalizeInteger(value)
	case schema.String:
		return normalizeString(value), nil
	case schema.Float:
		return normalizeFloat(value)
	case schema.Date, schema.Timestamp, schema.DateTime:
		return normalizeDateTime(value, fieldType)
	case schema.Object:
		return normalizeObject(value, fieldType.ObjectFields())
	default:
		return nil, errors.WrappedError(errInvalidValue, "%v, type: %T", value, value)
	}
}

func normalizeBoolean(value any) (bool, error) {
	switch value := value.(type) {
	case bool:
		return value, nil
	case string:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return false, errors.WrappedError(errInvalidBooleanValue, "string value: %s, error: %v", value, err)
		}

		return boolValue, nil
	default:
		return false, errors.WrappedError(errInvalidBooleanValue, "%v, type: %T", value, value)
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
			return 0, errors.WrappedError(errInvalidIntegerValue, "uint64 value: %d is greater than max int64 value", uintValue)
		}

		return int64(uintValue), nil
	case reflect.Float32, reflect.Float64:
		return int64(reflectVal.Float()), nil
	case reflect.String:
		value, err := normalizeNumberFromString(reflectVal.String())
		if err != nil {
			return 0, errors.WrappedError(errInvalidIntegerValue, "%s", err)
		}

		return int64(value), nil
	default:
		return 0, errors.WrappedError(errInvalidIntegerValue, "%v, type: %T", value, value)
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
			return 0, errors.WrappedError(errInvalidFloatValue, "%s", err)
		}

		return value, nil
	default:
		return 0, errors.WrappedError(errInvalidFloatValue, "%v, type: %T", value, value)
	}
}

func normalizeDateTime(value any, fieldType schema.FieldType) (any, error) {
	var (
		baseType  = fieldType.BaseType()
		timeValue time.Time
		err       error
	)

	switch value := value.(type) {
	case string:
		layout := fieldType.Layout()
		if layout == nil {
			return nil, errors.WrappedError(errInvalidValue, "layout is nil for %s", baseType)
		}

		timeValue, err = time.Parse(*layout, value)
		if err != nil {
			return nil, errors.WrappedError(errInvalidValue, "%s value: string %s cannot be formatted", baseType, value)
		}
	case time.Time:
		timeValue = value
	default:
		return nil, errors.WrappedError(errInvalidValue, "%v, type: %T", value, value)
	}

	layoutFormat, err := baseType.LayoutFormat()
	if err != nil {
		return nil, err
	}

	return timeValue.Format(layoutFormat), nil
}

func normalizeObject(value any, objectFields map[schema.Field]schema.FieldType) (any, error) {
	valueMap, ok := value.(map[string]any)
	if !ok {
		return nil, errors.WrappedError(errInvalidObjectValue, "%v", value)
	}

	normalizedObject := make(map[string]any, len(objectFields))
	for key, val := range valueMap {
		objectField, ok := objectFields[schema.Field(key)]
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
		if normalizedValue == nil {
			return nil
		}

		if slices.Contains(enumValues, normalizedValue) {
			return nil
		}

		return errors.WrappedError(errInvalidEnumValue, "value: %v", normalizedValue)
	}

	normalizedValueArray, ok := normalizedValue.([]any)
	// This shouldn't happen as previous normalization checks for this, but none the less we should check it
	if !ok {
		return errNormalizedValueIsNotAnArray
	}

	for _, value := range normalizedValueArray {
		if err := validateEnum(value, enumValues, false); err != nil {
			return err
		}
	}

	return nil
}
