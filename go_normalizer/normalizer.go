package gonormalizer

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"time"

	"diploma/rules"
)

func baseLayoutFormat(bt rules.BaseType) string {
	switch bt {
	case rules.Timestamp:
		return time.TimeOnly
	case rules.Date:
		return time.DateOnly
	case rules.DateTime:
		return time.DateTime
	default:
		return ""
	}
}

func layoutFormats(bt rules.BaseType) []string {
	switch bt {
	case rules.Timestamp:
		return timestampFormats()
	case rules.Date:
		return dateFormats()
	case rules.DateTime:
		return dateTimeFormats()
	default:
		return nil
	}
}

func timestampFormats() []string {
	return []string{
		time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano, time.TimeOnly,
	}
}

func dateFormats() []string {
	return []string{
		time.DateOnly,
	}
}

func dateTimeFormats() []string {
	return []string{
		time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850,
		time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano, time.DateTime,
	}
}

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
	for key, value := range data {
		fieldType, ok := n.Fields[rules.Field(key)]
		if !ok {
			delete(data, key)

			continue
		}

		if value == nil {
			continue
		}

		normalizedValue, err := normalize(value, fieldType)
		if err != nil {
			return nil, err
		}

		data[key] = normalizedValue
	}

	return data, nil
}

func normalize(value any, fieldType rules.FieldType) (normalizedValue any, err error) {
	if fieldType.IsArray() {
		normalizedValue, err = normalizeArray(value, fieldType)
	} else {
		normalizedValue, err = normalizeValue(value, fieldType)
	}

	if err != nil || len(fieldType.EnumValues()) == 0 {
		return
	}

	err = validateEnums(normalizedValue, fieldType.EnumValues(), fieldType.IsArray())

	return
}

func normalizeArray(value any, fieldType rules.FieldType) (any, error) {
	array, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("invalid array value: %v", value)
	}

	normalizedArray := make([]any, 0, len(array))
	for _, elem := range array {
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
		return fmt.Sprintf("%v", value), nil
	case rules.Float:
		return normalizeFloat(value)
	case rules.Date, rules.Timestamp, rules.DateTime:
		return normalizeDateTime(value, fieldType.BaseType())
	case rules.Object:
		return normalizeObject(value, fieldType)
	default:
		return nil, fmt.Errorf("invalid value: %v, type: %T", value, value)
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

		return nil, fmt.Errorf("invalid boolean value: string %s cannot be parsed to boolean", value)
	}

	return nil, fmt.Errorf("invalid boolean value: %v", value)
}

func normalizeInteger(value any) (any, error) {
	switch value := value.(type) {
	case int, int64:
		return value, nil
	case float64:
		return int64(value), nil
	case string:
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue, nil
		}

		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return int(math.Round(floatValue)), nil
		}

		return nil, fmt.Errorf("invalid integer value: string %s cannot be parsed to number", value)
	}

	return nil, fmt.Errorf("invalid integer value: %v, type: %T", value, value)
}

func normalizeFloat(value any) (any, error) {
	switch value := value.(type) {
	case int:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case float64:
		return value, nil
	case string:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue, nil
		}

		return nil, fmt.Errorf("invalid float value: string %s cannot be parsed to float", value)
	default:
		return nil, fmt.Errorf("invalid float value: %v, type: %T", value, value)
	}
}

func normalizeDateTime(value any, baseType rules.BaseType) (any, error) {
	stringValue, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("invalid %s value: %v", baseType, value)
	}

	layoutFormats := layoutFormats(baseType)
	for _, layoutFormat := range layoutFormats {
		parsedTime, err := time.Parse(layoutFormat, stringValue)
		if err == nil {
			return parsedTime.Format(baseLayoutFormat(baseType)), nil
		}
	}

	return nil, fmt.Errorf("invalid %s value: string %s cannot be formatted", baseType, stringValue)
}

func normalizeObject(value any, fieldType rules.FieldType) (any, error) {
	valueMap, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid object value: %v", value)
	}

	for key, val := range valueMap {
		objectField, ok := fieldType.ObjectFields()[rules.Field(key)]
		if !ok {
			delete(valueMap, key)

			continue
		}

		normalizedValue, err := normalize(val, objectField)
		if err != nil {
			return nil, err
		}

		valueMap[key] = normalizedValue
	}

	return valueMap, nil
}

func validateEnums(normalizedValue any, enumValues []any, isArray bool) error {
	if !isArray {
		if !slices.Contains(enumValues, normalizedValue) {
			return fmt.Errorf("invalid enum value %v", normalizedValue)
		}

		return nil
	}

	normalizedValueArray, ok := normalizedValue.([]any)
	// This shouldn't happen as previous normalization checks for this, but none the less we should check it
	if !ok {
		return errors.New("normalized value is not an array")
	}

	for _, value := range normalizedValueArray {
		if err := validateEnums(value, enumValues, false); err != nil {
			return err
		}
	}

	return nil
}
