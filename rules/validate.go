package rules

// ValidateNormalizer receives a generic function that creates a new normalizer
// and validates the normalizer returned by the function against the test cases.
func ValidateNormalizer(newNormalizer func(options ...NormalizerOption) (AbstractNormalizer, error)) error {
	normalizer, err := newNormalizer(testNormalizerOptions()...)
	if err != nil {
		return WrappedError(err, "failed to initialize normalizer")
	}

	return validateNormalizerTests(normalizer)
}

func validateNormalizerTests(normalizer AbstractNormalizer) error {
	for field, validateTest := range validateTests() {
		for _, inputOutput := range validateTest {
			inputMap := map[string]any{
				field: inputOutput.input,
			}

			output, err := normalizer.Normalize(inputMap)
			if err != nil {
				return WrappedError(err, "test %s failed", field)
			}

			if err := validateOutput(field, output, inputOutput.output); err != nil {
				return WrappedError(err, "test %s failed", field)
			}
		}
	}

	return nil
}

func validateOutput(field string, output map[string]any, expectedOutput any) error {
	if _, ok := output[field]; !ok && field != nonExistingTestField {
		return WrappedError(errExpectedOutputMismatch, "expected %v output, got nothing", expectedOutput)
	}

	if !compareValues(output[field], expectedOutput) {
		return WrappedError(errExpectedOutputMismatch, "expected %v output, got %v output", expectedOutput, output[field])
	}

	return nil
}
