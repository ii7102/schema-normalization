package jqNormalizer

import (
	"diploma/rules"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_JqNormalizer(t *testing.T) {
	err := rules.ValidateNormalizer(func(opts ...rules.NormalizerOption) (rules.AbstractNormalizer, error) {
		return NewNormalizer(opts...)
	})
	assert.NoError(t, err)
}
