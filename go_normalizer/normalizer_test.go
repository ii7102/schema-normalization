package gonormalizer_test

import (
	"testing"

	gonormalizer "github.com/ii7102/schema-normalization/go_normalizer"
	"github.com/ii7102/schema-normalization/schema"
	"github.com/ii7102/schema-normalization/validation"
	"github.com/stretchr/testify/assert"
)

func Test_GoNormalizer(t *testing.T) {
	t.Parallel()

	err := validation.ValidateNormalizer(func(opts ...schema.NormalizerOption) (schema.AbstractNormalizer, error) {
		return gonormalizer.NewNormalizer(opts...)
	})

	assert.NoError(t, err)
}
