package goNormalizer

import (
	"diploma/rules"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GoNormalizer(t *testing.T) {
	err := rules.ValidateNormalizer(func(opts ...rules.NormalizerOption) (rules.AbstractNormalizer, error) {
		return NewNormalizer(opts...)
	})
	assert.NoError(t, err)
}
