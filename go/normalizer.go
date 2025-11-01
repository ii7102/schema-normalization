package goNormalizer

import (
	"diploma/rules"
	"fmt"
	"math"
	"slices"
	"strconv"
)

var _ rules.AbstractNormalizer = (*normalizer)(nil)

type normalizer struct {
	rules.BaseNormalizer
}

func NewNormalizer(options ...rules.NormalizerOption) (*normalizer, error) {
	base, err := rules.NewBaseNormalizer(options...)
	if err != nil {
		return nil, err
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

		var normalizedValue any
		var err error
		if fieldType.IsArray {
			normalizedValue, err = normalizeArray(value, fieldType.BaseType)
		} else {
			normalizedValue, err = normalizeValue(value, fieldType.BaseType)
		}
		if err != nil {
			return nil, err
		}

		if len(fieldType.EnumValues) > 0 {
			if err = validateEnums(normalizedValue, fieldType.EnumValues, fieldType.IsArray); err != nil {
				return nil, err
			}
		}

		data[key] = normalizedValue
	}

	return data, nil
}

func normalizeArray(value any, fieldType rules.BaseType) (any, error) {
	array, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("invalid array value: %v", value)
	}

	for i, elem := range array {
		if elem == nil {
			continue
		}

		normalizedValue, err := normalizeValue(elem, fieldType)
		if err != nil {
			return nil, err
		}

		array[i] = normalizedValue
	}

	return array, nil
}

func normalizeValue(value any, fieldType rules.BaseType) (any, error) {
	if value == nil {
		return nil, nil
	}

	switch fieldType {
	case rules.Boolean:
		switch value := value.(type) {
		case bool:
			return value, nil
		case string:
			if boolValue, err := strconv.ParseBool(value); err == nil {
				return boolValue, nil
			}
		}
		return nil, fmt.Errorf("invalid boolean value: %v", value)
	case rules.Integer:
		switch value := value.(type) {
		case int, int64:
			return value, nil
		case float64:
			return int(math.Round(value)), nil
		case string:
			if intValue, err := strconv.Atoi(value); err == nil {
				return intValue, nil
			}
			if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
				return int(math.Round(floatValue)), nil
			}
		}
		return nil, fmt.Errorf("invalid integer value: %v, type: %T", value, value)
	case rules.String:
		return fmt.Sprintf("%v", value), nil
	case rules.Float:
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
		}
	}

	return value, nil
}

func validateEnums(normalizedValue any, enumValues []any, isArray bool) error {
	if !isArray {
		if !slices.Contains(enumValues, normalizedValue) {
			return fmt.Errorf("invalid enum value %v", normalizedValue)
		}

		return nil
	}

	normalizedValueArray, ok := normalizedValue.([]any)
	// This shouldn't happen as previous normalization checks for this, but nevertheless we should check it
	if !ok {
		return fmt.Errorf("normalized value is not an array")
	}

	for _, value := range normalizedValueArray {
		if err := validateEnums(value, enumValues, false); err != nil {
			return err
		}
	}

	return nil
}
