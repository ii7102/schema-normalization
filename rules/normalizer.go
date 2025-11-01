package rules

import "fmt"

type NormalizerOption func(AbstractNormalizer) error

func normalizerOptionError(err error) NormalizerOption {
	return func(n AbstractNormalizer) error {
		return err
	}
}

func withFields(fieldType FieldType, fields ...Field) NormalizerOption {
	return func(n AbstractNormalizer) error {
		for _, field := range fields {
			n.SetField(field, fieldType)
		}
		return nil
	}
}

func withArrayFields(fieldType FieldType, fields ...Field) NormalizerOption {
	return withFields(ArrayOf(fieldType), fields...)
}

func withEnumFields(fieldType FieldType, enumValues []any, fields ...Field) NormalizerOption {
	baseType := fieldType.BaseType
	if invalidEnumValue := validateEnumValues(baseType, enumValues); invalidEnumValue != nil {
		return normalizerOptionError(fmt.Errorf("enum value %v, isn't of %s type", invalidEnumValue, baseType.String()))
	}

	return withFields(EnumOf(fieldType, enumValues...), fields...)
}

func WithBooleanFields(fields ...Field) NormalizerOption {
	return withFields(BooleanType(), fields...)
}

func WithIntegerFields(fields ...Field) NormalizerOption {
	return withFields(IntegerType(), fields...)
}

func WithStringFields(fields ...Field) NormalizerOption {
	return withFields(StringType(), fields...)
}

func WithFloatFields(fields ...Field) NormalizerOption {
	return withFields(FloatType(), fields...)
}

func WithEnumOfBooleanFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(BooleanType(), enumValues, fields...)
}

func WithEnumOfIntegerFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(IntegerType(), enumValues, fields...)
}

func WithEnumOfStringFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(StringType(), enumValues, fields...)
}

func WithEnumOfFloatFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(FloatType(), enumValues, fields...)
}

func WithArrayOfBooleanFields(fields ...Field) NormalizerOption {
	return withArrayFields(BooleanType(), fields...)
}

func WithArrayOfIntegerFields(fields ...Field) NormalizerOption {
	return withArrayFields(IntegerType(), fields...)
}

func WithArrayOfStringFields(fields ...Field) NormalizerOption {
	return withArrayFields(StringType(), fields...)
}

func WithArrayOfFloatFields(fields ...Field) NormalizerOption {
	return withArrayFields(FloatType(), fields...)
}

func WithArrayOfEnumOfBooleanFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(ArrayOf(BooleanType()), enumValues, fields...)
}

func WithArrayOfEnumOfIntegerFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(ArrayOf(IntegerType()), enumValues, fields...)
}

func WithArrayOfEnumOfStringFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(ArrayOf(StringType()), enumValues, fields...)
}

func WithArrayOfEnumOfFloatFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(ArrayOf(FloatType()), enumValues, fields...)
}

type AbstractNormalizer interface {
	Normalize(map[string]any) (map[string]any, error)
	SetField(Field, FieldType)
	RemoveField(Field)
}

type BaseNormalizer struct {
	Fields map[Field]FieldType
}

func NewBaseNormalizer(options ...NormalizerOption) (BaseNormalizer, error) {
	n := BaseNormalizer{
		Fields: make(map[Field]FieldType),
	}

	for _, option := range options {
		if err := option(&n); err != nil {
			return BaseNormalizer{}, err
		}
	}

	return n, nil
}

func (n *BaseNormalizer) SetField(field Field, fieldType FieldType) {
	n.Fields[field] = fieldType
}

func (n *BaseNormalizer) RemoveField(field Field) {
	delete(n.Fields, field)
}

// Implement in the concrete normalizer
func (n *BaseNormalizer) Normalize(data map[string]any) (map[string]any, error) {
	return data, nil
}
