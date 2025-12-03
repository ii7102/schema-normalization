package validation

import (
	"github.com/ii7102/schema-normalization/errors"
	"github.com/ii7102/schema-normalization/schema"
)

// ValidateNormalizer receives a generic function that creates a new normalizer
// and validates the normalizer returned by the function against the test cases.
func ValidateNormalizer(newNormalizer func(_ ...schema.NormalizerOption) (schema.AbstractNormalizer, error)) error {
	normalizer, err := newNormalizer(normalizerTestOptions()...)
	if err != nil {
		return errors.WrappedError(err, "failed to initialize normalizer")
	}

	return validateNormalizerTests(normalizer)
}

func validateNormalizerTests(normalizer schema.AbstractNormalizer) error {
	for field, validateTest := range validateTests() {
		for _, inputOutput := range validateTest {
			inputMap := map[string]any{
				field: inputOutput.input,
			}

			output, err := normalizer.Normalize(inputMap)
			if err != nil {
				return errors.WrappedError(err, "test %s failed", field)
			}

			if err = validateOutput(field, output, inputOutput.output); err != nil {
				return errors.WrappedError(err, "test %s failed", field)
			}
		}
	}

	return nil
}

func validateOutput(field string, output map[string]any, expectedOutput any) error {
	if _, ok := output[field]; !ok && field != nonExistingTestField {
		return errors.WrappedError(errExpectedOutputMismatch, "expected %v, got nothing", expectedOutput)
	}

	if !compareValues(output[field], expectedOutput) {
		return errors.WrappedError(errExpectedOutputMismatch, "expected %v, got %v", expectedOutput, output[field])
	}

	return nil
}
