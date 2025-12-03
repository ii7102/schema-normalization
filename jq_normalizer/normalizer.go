package jqnormalizer

import (
	"github.com/ii7102/schema-normalization/errors"
	"github.com/ii7102/schema-normalization/schema"
	"github.com/itchyny/gojq"
)

var _ schema.AbstractNormalizer = (*normalizer)(nil)

type normalizer struct {
	*schema.BaseNormalizer

	cachedCode *gojq.Code
}

// NewNormalizer creates a new JQ normalizer with the given options and caches the compiled JQ code.
func NewNormalizer(opts ...schema.NormalizerOption) (*normalizer, error) {
	base, err := schema.NewBaseNormalizer(opts...)
	if err != nil {
		return nil, errors.WrappedError(err, "failed to create base normalizer")
	}

	return &normalizer{
		BaseNormalizer: base,
		cachedCode:     compileJqCode(jqFilter(base.Fields(), 0)),
	}, nil
}

func (n *normalizer) SetField(field schema.Field, fieldType schema.FieldType) {
	n.BaseNormalizer.SetField(field, fieldType)
	n.CacheCompiledJqCode()
}

func (n *normalizer) RemoveField(field schema.Field) {
	if !n.HasField(field) {
		return
	}

	n.BaseNormalizer.RemoveField(field)
	n.CacheCompiledJqCode()
}

func (n *normalizer) JqFilter() string {
	return jqFilter(n.Fields(), 0)
}

func (n *normalizer) JqBatchFilter() string {
	return jqBatchFilter(n.Fields())
}

func (n *normalizer) CacheCompiledJqCode() {
	n.cachedCode = compileJqCode(n.JqFilter())
}

func (n *normalizer) Normalize(data map[string]any) (map[string]any, error) {
	if n.cachedCode == nil {
		return nil, errJqQueryCompilationNotCached
	}

	iter := n.cachedCode.Run(data)

	result, ok := iter.Next()
	if !ok {
		return nil, errNoResultFromJqQuery
	}

	switch result := result.(type) {
	case error:
		return nil, result
	case map[string]any:
		return result, nil
	default:
		return nil, errors.WrappedError(errUnexpectedResultType, "unexpected result type: %T", result)
	}
}

// NormalizeBatch normalizes a batch of data using the cached JQ code.
// Batching in jq reduces memory allocations by ~50% but provides minimal throughput improvement (~3%),
// indicating that the jq evaluation itself—not per-call overhead—is the primary performance bottleneck.
func (n *normalizer) NormalizeBatch(data []any) ([]any, error) {
	jqQuery, err := gojq.Parse(n.JqBatchFilter())
	if err != nil {
		return nil, errors.WrappedError(err, "failed to parse JQ batch filter")
	}

	iter := jqQuery.Run(data)

	result, ok := iter.Next()
	if !ok {
		return nil, errNoResultFromJqQuery
	}

	switch result := result.(type) {
	case error:
		return nil, result
	case []any:
		return result, nil
	default:
		return nil, errors.WrappedError(errUnexpectedResultType, "unexpected result type: %T", result)
	}
}
