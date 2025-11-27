package rules

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
			return nil, WrappedError(err, "failed to apply normalizer option")
		}
	}

	return normalizer, nil
}

// SetField sets the given field and its type to the normalizer.
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

// NormalizeBatch is a placeholder for the concrete normalizer implementation.
func (*BaseNormalizer) NormalizeBatch(data []any) ([]any, error) {
	return data, nil
}
