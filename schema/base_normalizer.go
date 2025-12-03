package schema

import "github.com/ii7102/schema-normalization/errors"

// BaseNormalizer is the base generic normalizer that implements generic SetField and RemoveField methods.
// To be compatible with the AbstractNormalizer interface, it implements a placeholder for the Normalize method.
// The concrete normalizers that extend the BaseNormalizer should implement the Normalize method and
// optionally override the SetField and RemoveField methods.
type BaseNormalizer struct {
	fields map[Field]FieldType
}

// NewBaseNormalizer creates a new base normalizer with the given options.
func NewBaseNormalizer(opts ...NormalizerOption) (*BaseNormalizer, error) {
	normalizer := &BaseNormalizer{
		fields: make(map[Field]FieldType),
	}

	for _, opt := range opts {
		if err := opt(normalizer); err != nil {
			return nil, errors.WrappedError(err, "failed to apply normalizer option")
		}
	}

	return normalizer, nil
}

// SetField sets the given field and its type to the normalizer.
func (n *BaseNormalizer) SetField(field Field, fieldType FieldType) {
	n.fields[field] = fieldType
}

// RemoveField removes the given field from the normalizer.
func (n *BaseNormalizer) RemoveField(field Field) {
	delete(n.fields, field)
}

// Normalize is a placeholder for the concrete normalizer implementation.
func (*BaseNormalizer) Normalize(data map[string]any) (map[string]any, error) {
	return data, nil
}

// NormalizeBatch is a placeholder for the concrete normalizer implementation.
func (*BaseNormalizer) NormalizeBatch(data []any) ([]any, error) {
	return data, nil
}

// Fields returns the fields of the normalizer.
func (n *BaseNormalizer) Fields() map[Field]FieldType {
	return n.fields
}

// FieldType returns the field type for the given field.
func (n *BaseNormalizer) FieldType(field Field) FieldType {
	return n.fields[field]
}

// HasField returns true if the given field is in the normalizer and false otherwise.
func (n *BaseNormalizer) HasField(field Field) bool {
	_, ok := n.fields[field]

	return ok
}
