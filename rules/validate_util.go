package rules

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
)

const (
	nonExistingTestField = "nonExistingTestField"
)

func stringEnumValues() []any {
	return []any{"test_string", "0", "1", "-1.1", "1.1", "true", "false", "", nil}
}

func floatEnumValues() []any {
	return []any{0.0, 1.0, -1.0, -1.1, 1.1, 0.0, -1.0, 1.0, nil}
}

func testNormalizerOptions() []NormalizerOption {
	return []NormalizerOption{
		WithBooleanFields(BooleanType().Field()),
		WithIntegerFields(IntegerType().Field()),
		WithStringFields(StringType().Field()),
		WithFloatFields(FloatType().Field()),
		WithEnumOfStringFields(stringEnumValues(), EnumOf(StringType(), nil).Field()),
		WithArrayOfBooleanFields(ArrayOf(BooleanType()).Field()),
		WithArrayOfIntegerFields(ArrayOf(IntegerType()).Field()),
		WithArrayOfStringFields(ArrayOf(StringType()).Field()),
		WithArrayOfFloatFields(ArrayOf(FloatType()).Field()),
		WithArrayOfEnumOfFloatFields(floatEnumValues(), EnumOf(FloatType(), nil).Field()),
	}
}

type inputOutput struct {
	input, output any
}

func validateTests() map[string][]inputOutput {
	randInt := rand.Int()
	randInt63 := rand.Int63()
	randFloat64 := rand.Float64()
	randBoolean := randomBoolean()
	randString := randomString()

	return map[string][]inputOutput{
		nonExistingTestField: {
			{input: true},
			{input: 0},
			{input: 0.123456789},
			{input: "test_string"},
			{input: randInt},
			{input: randInt63},
			{input: randFloat64},
			{input: randBoolean},
			{input: randString},
			{input: nil},
		},
		BooleanType().String(): {
			{input: true, output: true},
			{input: false, output: false},
			{input: "true", output: true},
			{input: "false", output: false},
			{input: nil, output: nil},
		},
		IntegerType().String(): {
			{input: 0, output: 0},
			{input: 1, output: 1},
			{input: 1.0, output: 1},
			{input: 1.1, output: 1},
			{input: 1.9, output: 2},
			{input: -1.0, output: -1},
			{input: -1.1, output: -1},
			{input: -1.9, output: -2},
			{input: 123456789, output: 123456789},
			{input: 0.123456789, output: 0},
			{input: "3", output: 3},
			{input: randInt, output: randInt},
			{input: randInt63, output: randInt63},
			{input: randFloat64, output: int(math.Round(randFloat64))},
			{input: nil, output: nil},
		},
		StringType().String(): {
			{input: "test_string", output: "test_string"},
			{input: "3", output: "3"},
			{input: "", output: ""},
			{input: true, output: "true"},
			{input: false, output: "false"},
			{input: 0, output: "0"},
			{input: 1.1, output: "1.1"},
			{input: 1.0, output: "1"},
			{input: -1.1, output: "-1.1"},
			{input: 123456789, output: "123456789"},
			{input: randString, output: randString},
			{input: randInt, output: fmt.Sprintf("%d", randInt)},
			{input: randInt63, output: fmt.Sprintf("%d", randInt63)},
			{input: randFloat64, output: fmt.Sprintf("%v", randFloat64)},
			{input: nil, output: nil},
		},
		FloatType().String(): {
			{input: 0, output: 0.0},
			{input: 1, output: 1.0},
			{input: 1.0, output: 1.0},
			{input: 1.1, output: 1.1},
			{input: 1.9, output: 1.9},
			{input: -1.0, output: -1.0},
			{input: -1.1, output: -1.1},
			{input: -1.9, output: -1.9},
			{input: 123456789, output: 123456789.0},
			{input: 0.123456789, output: 0.123456789},
			{input: "3", output: 3.0},
			{input: randInt, output: float64(randInt)},
			{input: randInt63, output: float64(randInt63)},
			{input: randFloat64, output: randFloat64},
			{input: nil, output: nil},
		},
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString() string {
	s := make([]rune, rand.Intn(100))
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func randomBoolean() bool {
	return rand.Intn(2) == 0
}

func compareValues(v1, v2 any) bool {
	if v1 == nil && v2 == nil {
		return true
	}
	if v1 == nil || v2 == nil {
		return false
	}

	val1 := reflect.ValueOf(v1)
	val2 := reflect.ValueOf(v2)

	if val1.Kind() == reflect.Slice && val2.Kind() == reflect.Slice {
		return compareSlices(val1, val2)
	}

	return v1 == v2
}

func compareSlices(slice1, slice2 reflect.Value) bool {
	if slice1.Len() != slice2.Len() {
		return false
	}

	for i := 0; i < slice1.Len(); i++ {
		if !compareValues(slice1.Index(i).Interface(), slice2.Index(i).Interface()) {
			return false
		}
	}

	return true
}
