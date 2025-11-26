package jqnormalizer_test

import (
	"testing"

	jqnormalizer "github.com/ii7102/schema-normalization/jq_normalizer"
	"github.com/ii7102/schema-normalization/rules"
	"github.com/stretchr/testify/assert"
)

func Test_JqNormalizer(t *testing.T) {
	t.Parallel()

	err := rules.ValidateNormalizer(func(opts ...rules.NormalizerOption) (rules.AbstractNormalizer, error) {
		return jqnormalizer.NewNormalizer(opts...)
	})

	assert.NoError(t, err)
}
