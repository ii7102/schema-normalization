package jqNormalizer

import (
	"fmt"
	"sync"

	"diploma/rules"

	"github.com/itchyny/gojq"
)

var _ rules.AbstractNormalizer = (*normalizer)(nil)

type normalizer struct {
	rules.BaseNormalizer
	cachedCode     *gojq.Code
	fieldCodeCache map[string]*gojq.Code
	cacheMutex     sync.RWMutex
}

func NewNormalizer(options ...rules.NormalizerOption) (*normalizer, error) {
	base, err := rules.NewBaseNormalizer(options...)
	if err != nil {
		return nil, err
	}

	n := &normalizer{
		BaseNormalizer: base,
		fieldCodeCache: make(map[string]*gojq.Code),
	}
	n.CacheCompiledJqCode()
	n.precompileFieldRules()
	return n, nil
}

func (n *normalizer) SetField(field rules.Field, fieldType rules.FieldType) {
	n.BaseNormalizer.SetField(field, fieldType)
	n.CacheCompiledJqCode()
	n.precompileFieldRules()
}

func (n *normalizer) RemoveField(field rules.Field) {
	if _, ok := n.Fields[field]; !ok {
		return
	}
	n.BaseNormalizer.RemoveField(field)
	n.CacheCompiledJqCode()
	n.precompileFieldRules()
}

func (n *normalizer) JqFilter() string {
	return jqFilter(n.Fields)
}

func (n *normalizer) CacheCompiledJqCode() {
	n.cachedCode = compileJqCode(n.JqFilter())
}

func (n *normalizer) precompileFieldRules() {
	n.cacheMutex.Lock()
	defer n.cacheMutex.Unlock()

	jqRules := jqRules()
	for fieldTypeStr, jqRule := range jqRules {
		if n.fieldCodeCache[fieldTypeStr] == nil {
			n.fieldCodeCache[fieldTypeStr] = compileJqCode(jqRule)
		}
	}
}

func (n *normalizer) Normalize(data map[string]any) (map[string]any, error) {
	if len(n.Fields) == 0 {
		return data, nil
	}

	// Use optimized field-by-field processing with pre-compiled jq rules
	return n.normalizeFieldByField(data)
}

func (n *normalizer) normalizeFieldByField(data map[string]any) (map[string]any, error) {
	result := make(map[string]any, len(n.Fields))

	n.cacheMutex.RLock()
	defer n.cacheMutex.RUnlock()

	for field, fieldType := range n.Fields {
		fieldStr := string(field)
		value, exists := data[fieldStr]
		if !exists {
			continue
		}

		if value == nil {
			result[fieldStr] = nil
			continue
		}

		// Use pre-compiled jq code for this field type
		fieldTypeStr := fieldType.String()
		compiledCode, ok := n.fieldCodeCache[fieldTypeStr]
		if !ok {
			// Fallback to original method if no cached code
			return n.normalizeFallback(data)
		}

		// Apply jq transformation
		iter := compiledCode.Run(value)
		normalizedValue, ok := iter.Next()
		if !ok {
			result[fieldStr] = nil
			continue
		}

		switch normalizedValue := normalizedValue.(type) {
		case error:
			return nil, normalizedValue
		default:
			// Validate enums if specified
			if len(fieldType.EnumValues) > 0 {
				if err := n.validateEnums(normalizedValue, fieldType.EnumValues, fieldType.IsArray); err != nil {
					return nil, err
				}
			}
			result[fieldStr] = normalizedValue
		}
	}

	return result, nil
}

// Simple enum validation without external dependencies
func (n *normalizer) validateEnums(normalizedValue any, enumValues []any, isArray bool) error {
	if !isArray {
		for _, enumValue := range enumValues {
			if normalizedValue == enumValue {
				return nil
			}
		}
		return fmt.Errorf("invalid enum value %v", normalizedValue)
	}

	normalizedValueArray, ok := normalizedValue.([]any)
	if !ok {
		return fmt.Errorf("normalized value is not an array")
	}

	for _, value := range normalizedValueArray {
		if err := n.validateEnums(value, enumValues, false); err != nil {
			return err
		}
	}

	return nil
}

func (n *normalizer) normalizeFallback(data map[string]any) (map[string]any, error) {
	if n.cachedCode == nil {
		return nil, fmt.Errorf("jq query compilation not cached")
	}

	iter := n.cachedCode.Run(data)
	result, ok := iter.Next()
	if !ok {
		return nil, nil
	}

	switch result := result.(type) {
	case error:
		return nil, result
	case map[string]any:
		// Remove fields that weren't in original data
		for key := range result {
			if _, ok := data[key]; !ok {
				delete(result, key)
			}
		}
		return result, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}
}
