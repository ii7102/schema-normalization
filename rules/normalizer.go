package rules

// NormalizerOption is a function that sets fields of some type to the normalizer.
type NormalizerOption func(AbstractNormalizer) error

func normalizerOptionError(err error) NormalizerOption {
	return func(_ AbstractNormalizer) error {
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
	enum, err := EnumOf(fieldType, enumValues...)
	if err != nil {
		return normalizerOptionError(err)
	}

	return withFields(enum, fields...)
}

// WithBooleanFields sets the given fields as a boolean type.
func WithBooleanFields(fields ...Field) NormalizerOption {
	return withFields(BooleanType(), fields...)
}

// WithIntegerFields sets the given fields as a integer type.
func WithIntegerFields(fields ...Field) NormalizerOption {
	return withFields(IntegerType(), fields...)
}

// WithStringFields sets the given fields as a string type.
func WithStringFields(fields ...Field) NormalizerOption {
	return withFields(StringType(), fields...)
}

// WithFloatFields sets the given fields as a float type.
func WithFloatFields(fields ...Field) NormalizerOption {
	return withFields(FloatType(), fields...)
}

// WithDateFields sets the given fields as a date type.
func WithDateFields(fields ...Field) NormalizerOption {
	return withFields(DateType(), fields...)
}

// WithTimestampFields sets the given fields as a timestamp type.
func WithTimestampFields(fields ...Field) NormalizerOption {
	return withFields(TimestampType(), fields...)
}

// WithDateTimeFields sets the given fields as a dateTime type.
func WithDateTimeFields(fields ...Field) NormalizerOption {
	return withFields(DateTimeType(), fields...)
}

// WithObjectField sets the given fields as a object type with the given object fields.
func WithObjectField(objectFields map[Field]FieldType, fields ...Field) NormalizerOption {
	objectType, err := ObjectType(objectFields)
	if err != nil {
		return normalizerOptionError(err)
	}

	return withFields(objectType, fields...)
}

// WithEnumOfBooleanFields sets the given fields as a enum of boolean type with the given enum values.
func WithEnumOfBooleanFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(BooleanType(), enumValues, fields...)
}

// WithEnumOfIntegerFields sets the given fields as a enum of integer type with the given enum values.
func WithEnumOfIntegerFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(IntegerType(), enumValues, fields...)
}

// WithEnumOfStringFields sets the given fields as a enum of string type with the given enum values.
func WithEnumOfStringFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(StringType(), enumValues, fields...)
}

// WithEnumOfFloatFields sets the given fields as a enum of float type with the given enum values.
func WithEnumOfFloatFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(FloatType(), enumValues, fields...)
}

// WithArrayOfBooleanFields sets the given fields as a array of boolean type.
func WithArrayOfBooleanFields(fields ...Field) NormalizerOption {
	return withArrayFields(BooleanType(), fields...)
}

// WithArrayOfIntegerFields sets the given fields as a array of integer type.
func WithArrayOfIntegerFields(fields ...Field) NormalizerOption {
	return withArrayFields(IntegerType(), fields...)
}

// WithArrayOfStringFields sets the given fields as a array of string type.
func WithArrayOfStringFields(fields ...Field) NormalizerOption {
	return withArrayFields(StringType(), fields...)
}

// WithArrayOfFloatFields sets the given fields as a array of float type.
func WithArrayOfFloatFields(fields ...Field) NormalizerOption {
	return withArrayFields(FloatType(), fields...)
}

// WithArrayOfEnumOfBooleanFields sets the given fields as a array of enum of boolean type with the given enum values.
func WithArrayOfEnumOfBooleanFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(ArrayOf(BooleanType()), enumValues, fields...)
}

// WithArrayOfEnumOfIntegerFields sets the given fields as a array of enum of integer type with the given enum values.
func WithArrayOfEnumOfIntegerFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(ArrayOf(IntegerType()), enumValues, fields...)
}

// WithArrayOfEnumOfStringFields sets the given fields as a array of enum of string type with the given enum values.
func WithArrayOfEnumOfStringFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(ArrayOf(StringType()), enumValues, fields...)
}

// WithArrayOfEnumOfFloatFields sets the given fields as a array of enum of float type with the given enum values.
func WithArrayOfEnumOfFloatFields(enumValues []any, fields ...Field) NormalizerOption {
	return withEnumFields(ArrayOf(FloatType()), enumValues, fields...)
}

// AbstractNormalizer is the abstract interface that all normalizers must implement.
// It defines the methods that all normalizers must implement.
type AbstractNormalizer interface {
	Normalize(data map[string]any) (map[string]any, error)
	SetField(field Field, fieldType FieldType)
	RemoveField(field Field)
}

// BaseNormalizer is the base generic normalizer that implements generic SetField and RemoveField methods.
// To be compatible with the AbstractNormalizer interface, it implements a placeholder for the Normalize method.
// The concrete normalizers that extend the BaseNormalizer should implement the Normalize method and
// optionally override the SetField and RemoveField methods.
type BaseNormalizer struct {
	Fields map[Field]FieldType
}

// NewBaseNormalizer creates a new base normalizer with the given options.
func NewBaseNormalizer(options ...NormalizerOption) (*BaseNormalizer, error) {
	normalizer := &BaseNormalizer{
		Fields: make(map[Field]FieldType),
	}

	for _, option := range options {
		if err := option(normalizer); err != nil {
			return nil, err
		}
	}

	return normalizer, nil
}

// SetField sets the given field and it's type to the normalizer.
func (n *BaseNormalizer) SetField(field Field, fieldType FieldType) {
	n.Fields[field] = fieldType
}

// RemoveField removes the given field from the normalizer.
func (n *BaseNormalizer) RemoveField(field Field) {
	delete(n.Fields, field)
}

// Normalize is a placeholder for the concrete normalizer implementation.
func (*BaseNormalizer) Normalize(data map[string]any) (map[string]any, error) {
	return data, nil
}
