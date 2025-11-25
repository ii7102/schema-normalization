package gonormalizer_test

import (
	"testing"

	gonormalizer "github.com/ii7102/schema-normalization/go_normalizer"
	"github.com/ii7102/schema-normalization/rules"
	"github.com/stretchr/testify/assert"
)

func Test_GoNormalizer(t *testing.T) {
	t.Parallel()

	err := rules.ValidateNormalizer(func(opts ...rules.NormalizerOption) (rules.AbstractNormalizer, error) {
		return gonormalizer.NewNormalizer(opts...)
	})
	assert.NoError(t, err)
}
