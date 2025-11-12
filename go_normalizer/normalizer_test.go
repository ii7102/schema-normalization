package gonormalizer_test

import (
	"testing"

	gonormalizer "diploma/go_normalizer"
	"diploma/rules"
	"github.com/stretchr/testify/assert"
)

func Test_GoNormalizer(t *testing.T) {
	t.Parallel()

	err := rules.ValidateNormalizer(func(opts ...rules.NormalizerOption) (rules.AbstractNormalizer, error) {
		return gonormalizer.NewNormalizer(opts...)
	})
	assert.NoError(t, err)
}
