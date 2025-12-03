package jqnormalizer_test

import (
	"testing"

	jqnormalizer "github.com/ii7102/schema-normalization/jq_normalizer"
	"github.com/ii7102/schema-normalization/schema"
	"github.com/ii7102/schema-normalization/validation"
	"github.com/stretchr/testify/assert"
)

func Test_JqNormalizer(t *testing.T) {
	t.Parallel()

	err := validation.ValidateNormalizer(func(opts ...schema.NormalizerOption) (schema.AbstractNormalizer, error) {
		return jqnormalizer.NewNormalizer(opts...)
	})

	assert.NoError(t, err)
}
