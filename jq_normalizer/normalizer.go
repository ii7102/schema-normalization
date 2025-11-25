package jqnormalizer

import (
	"github.com/ii7102/schema-normalization/rules"
	"github.com/itchyny/gojq"
)

var _ rules.AbstractNormalizer = (*normalizer)(nil)

type normalizer struct {
	*rules.BaseNormalizer

	cachedCode *gojq.Code
}

// NewNormalizer creates a new JQ normalizer with the given options and caches the compiled JQ code.
func NewNormalizer(options ...rules.NormalizerOption) (*normalizer, error) {
	base, err := rules.NewBaseNormalizer(options...)
	if err != nil {
		return nil, rules.WrappedError(err, "failed to create base normalizer")
	}

	n := &normalizer{
		BaseNormalizer: base,
		cachedCode:     compileJqCode(jqFilter(base.Fields)),
	}

	return n, nil
}

func (n *normalizer) SetField(field rules.Field, fieldType rules.FieldType) {
	n.BaseNormalizer.SetField(field, fieldType)
	n.CacheCompiledJqCode()
}

func (n *normalizer) RemoveField(field rules.Field) {
	if _, ok := n.Fields[field]; !ok {
		return
	}

	n.BaseNormalizer.RemoveField(field)
	n.CacheCompiledJqCode()
}

func (n *normalizer) JqFilter() string {
	return jqFilter(n.Fields)
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
		return nil, rules.WrappedError(errUnexpectedResultType, "unexpected result type: %T", result)
	}
}
