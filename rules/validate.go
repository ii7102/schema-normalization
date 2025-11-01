package rules

import (
	"fmt"
)

func ValidateNormalizer(newNormalizer func(options ...NormalizerOption) (AbstractNormalizer, error)) error {
	normalizer, err := newNormalizer(testNormalizerOptions()...)
	if err != nil {
		return fmt.Errorf("failed to initialize normalizer, error: %v", err)
	}

	for field, validateTest := range validateTests() {
		for _, inputOutput := range validateTest {
			inputMap := map[string]any{
				field: inputOutput.input,
			}

			output, err := normalizer.Normalize(inputMap)

			if err != nil {
				return fmt.Errorf("test %s failed: %v", field, err)
			}

			if _, ok := output[field]; !ok && field != nonExistingTestField {
				return fmt.Errorf("test %s failed: expected %v output, got nothing", field, inputOutput.output)
			}

			if !compareValues(output[field], inputOutput.output) {
				return fmt.Errorf("test %s failed: expected %v output, got %v output", field, inputOutput.output, output[field])
			}
		}
	}

	return nil
}
