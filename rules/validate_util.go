package rules

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
)

const (
	nonExistingTestField = "nonExistingTestField"

	maxStringLength = 100
	maxBoolean      = 2
)

func stringEnumValues() []any {
	return []any{"test_string", "0", "1", "-1.1", "1.1", "true", "false", ""}
}

func floatEnumValues() []any {
	return []any{0.0, 1.0, -1.0, -1.1, 1.1, 0.0}
}

func objectFields() map[Field]FieldType {
	return map[Field]FieldType{
		BooleanType().Field(): BooleanType(),
		IntegerType().Field(): IntegerType(),
		StringType().Field():  StringType(),
		FloatType().Field():   FloatType(),
	}
}

func testNormalizerOptions() []NormalizerOption {
	stringEnum, _ := EnumOf(StringType(), stringEnumValues()...)
	floatEnum, _ := EnumOf(FloatType(), floatEnumValues()...)
	objectType, _ := ObjectType(objectFields())

	return []NormalizerOption{
		WithBooleanFields(BooleanType().Field()),
		WithIntegerFields(IntegerType().Field()),
		WithStringFields(StringType().Field()),
		WithFloatFields(FloatType().Field()),
		WithEnumOfStringFields(stringEnumValues(), stringEnum.Field()),
		WithArrayOfBooleanFields(ArrayOf(BooleanType()).Field()),
		WithArrayOfIntegerFields(ArrayOf(IntegerType()).Field()),
		WithArrayOfStringFields(ArrayOf(StringType()).Field()),
		WithArrayOfFloatFields(ArrayOf(FloatType()).Field()),
		WithArrayOfEnumOfFloatFields(floatEnumValues(), ArrayOf(floatEnum).Field()),
		WithObjectField(objectFields(), objectType.Field()),
		WithDateFields(DateType().Field()),
		WithTimestampFields(TimestampType().Field()),
		WithDateTimeFields(DateTimeType().Field()),
	}
}

type inputOutput struct {
	input, output any
}

type randomTestValues struct {
	int     int
	int63   int64
	float63 float64
	boolean bool
	string  string
}

func generateRandomTestValues() randomTestValues {
	return randomTestValues{
		int:     rand.Int(),
		int63:   rand.Int63(),
		float63: rand.Float64(),
		boolean: randomBoolean(),
		string:  randomString(),
	}
}

func nonExistingFieldTests(rnd randomTestValues) []inputOutput {
	return []inputOutput{
		{input: true, output: nil},
		{input: 0, output: nil},
		{input: 0.123456789, output: nil},
		{input: "test_string", output: nil},
		{input: rnd.int, output: nil},
		{input: rnd.int63, output: nil},
		{input: rnd.float63, output: nil},
		{input: rnd.boolean, output: nil},
		{input: rnd.string, output: nil},
		{input: nil, output: nil},
	}
}

func booleanTypeTests() []inputOutput {
	return []inputOutput{
		{input: true, output: true},
		{input: false, output: false},
		{input: "true", output: true},
		{input: "false", output: false},
		{input: nil, output: nil},
	}
}

func integerTypeTests(rnd randomTestValues) []inputOutput {
	return []inputOutput{
		{input: 0, output: 0},
		{input: 1, output: 1},
		{input: 1.0, output: 1},
		{input: 1.1, output: 1},
		{input: 1.9, output: 1},
		{input: -1.0, output: -1},
		{input: -1.1, output: -1},
		{input: -1.9, output: -1},
		{input: 123456789, output: 123456789},
		{input: 0.123456789, output: 0},
		{input: "3", output: 3},
		{input: rnd.int, output: rnd.int},
		{input: rnd.int63, output: rnd.int63},
		{input: rnd.float63, output: int64(rnd.float63)},
		{input: nil, output: nil},
	}
}

func stringTypeTests(rnd randomTestValues) []inputOutput {
	return []inputOutput{
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
		{input: rnd.string, output: rnd.string},
		{input: rnd.int, output: strconv.Itoa(rnd.int)},
		{input: rnd.int63, output: strconv.FormatInt(rnd.int63, 10)},
		{input: rnd.float63, output: fmt.Sprintf("%v", rnd.float63)},
		{input: nil, output: nil},
	}
}

func floatTypeTests(rnd randomTestValues) []inputOutput {
	return []inputOutput{
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
		{input: rnd.int, output: float64(rnd.int)},
		{input: rnd.int63, output: float64(rnd.int63)},
		{input: rnd.float63, output: rnd.float63},
		{input: nil, output: nil},
	}
}

func stringEnumTests() []inputOutput {
	return []inputOutput{
		{input: "test_string", output: "test_string"},
		{input: "0", output: "0"},
		{input: "1", output: "1"},
		{input: "-1.1", output: "-1.1"},
		{input: "1.1", output: "1.1"},
		{input: "true", output: "true"},
		{input: "false", output: "false"},
		{input: "", output: ""},
		{input: nil, output: nil},
		{input: 0, output: "0"},
		{input: 1, output: "1"},
		{input: -1.1, output: "-1.1"},
		{input: 1.1, output: "1.1"},
		{input: true, output: "true"},
		{input: false, output: "false"},
	}
}

func arrayOfBooleanTypeTests() []inputOutput {
	return []inputOutput{
		{
			input:  []any{true, false, true, false},
			output: []any{true, false, true, false},
		},
		{
			input:  []any{"true", "false", "true", "false"},
			output: []any{true, false, true, false},
		},
		{
			input:  []any{nil, nil, nil, nil},
			output: []any{nil, nil, nil, nil},
		},
		{
			input:  []any{true, false, "true", "false", nil, nil},
			output: []any{true, false, true, false, nil, nil},
		},
		{
			input:  []any{},
			output: []any{},
		},
	}
}

func arrayOfIntegerTypeTests(rnd randomTestValues) []inputOutput {
	return []inputOutput{
		{
			input:  []any{0, 1, 1.0, 1.1, 1.9},
			output: []any{0, 1, 1, 1, 1},
		},
		{
			input:  []any{0, -1, -1.0, -1.1, -1.9},
			output: []any{0, -1, -1, -1, -1},
		},
		{
			input:  []any{123456789, 0.123456789},
			output: []any{123456789, 0},
		},
		{
			input:  []any{"0", "1", "2.0", "-1", "-2.0"},
			output: []any{0, 1, 2, -1, -2},
		},
		{
			input:  []any{nil, nil, nil, nil},
			output: []any{nil, nil, nil, nil},
		},
		{
			input:  []any{0, 1, 1.0, "-1", "-1.0", nil, nil},
			output: []any{0, 1, 1, -1, -1, nil, nil},
		},
		{
			input:  []any{},
			output: []any{},
		},
		{
			input:  []any{rnd.int, rnd.int63, rnd.float63},
			output: []any{rnd.int, rnd.int63, int64(rnd.float63)},
		},
	}
}

func arrayOfStringTypeTests(rnd randomTestValues) []inputOutput {
	return []inputOutput{
		{
			input:  []any{"test_string", "3", "3.0", "", "true", "false"},
			output: []any{"test_string", "3", "3.0", "", "true", "false"},
		},
		{
			input:  []any{true, true, false, false},
			output: []any{"true", "true", "false", "false"},
		},
		{
			input:  []any{0, 1, 1.0, 1.1, -1.1, 123456789},
			output: []any{"0", "1", "1", "1.1", "-1.1", "123456789"},
		},
		{
			input:  []any{nil, nil, nil, nil},
			output: []any{nil, nil, nil, nil},
		},
		{
			input:  []any{"test_string", "3", "3.0", "", "true", "false", 3, 3.0, nil, true, false, nil},
			output: []any{"test_string", "3", "3.0", "", "true", "false", "3", "3", nil, "true", "false", nil},
		},
		{
			input:  []any{},
			output: []any{},
		},
		{
			input:  []any{rnd.string, rnd.int, rnd.int63, rnd.float63},
			output: []any{rnd.string, strconv.Itoa(rnd.int), strconv.FormatInt(rnd.int63, 10), fmt.Sprintf("%v", rnd.float63)},
		},
	}
}

func arrayOfFloatTypeTests(rnd randomTestValues) []inputOutput {
	return []inputOutput{
		{
			input:  []any{0, 1, 1.0, 1.1, 1.9},
			output: []any{0.0, 1.0, 1.0, 1.1, 1.9},
		},
		{
			input:  []any{0, -1, -1.0, -1.1, -1.9},
			output: []any{0.0, -1.0, -1.0, -1.1, -1.9},
		},
		{
			input:  []any{123456789, 0.123456789},
			output: []any{123456789.0, 0.123456789},
		},
		{
			input:  []any{"0", "1", "2.0", "-1", "-2.0"},
			output: []any{0.0, 1.0, 2.0, -1.0, -2.0},
		},
		{
			input:  []any{nil, nil, nil, nil},
			output: []any{nil, nil, nil, nil},
		},
		{
			input:  []any{0, 1, 1.0, "-1", "-1.0", nil, nil},
			output: []any{0.0, 1.0, 1.0, -1.0, -1.0, nil, nil},
		},
		{
			input:  []any{},
			output: []any{},
		},
		{
			input:  []any{rnd.int, rnd.int63, rnd.float63},
			output: []any{float64(rnd.int), float64(rnd.int63), rnd.float63},
		},
	}
}

func arrayOfFloatEnumTests() []inputOutput {
	return []inputOutput{
		{
			input:  []any{0.0, 1.0, -1.0, -1.1, 1.1, 0.0, 1.0, -1.0, -1.1, 1.1},
			output: []any{0.0, 1.0, -1.0, -1.1, 1.1, 0.0, 1.0, -1.0, -1.1, 1.1},
		},
		{
			input:  []any{"0.0", "1.0", "-1.0", "-1.1", "1.1"},
			output: []any{0.0, 1.0, -1.0, -1.1, 1.1},
		},
		{
			input:  []any{0, 1, -1},
			output: []any{0.0, 1.0, -1.0},
		},
		{
			input:  []any{nil, nil, nil, nil},
			output: []any{nil, nil, nil, nil},
		},
		{
			input:  []any{0, 1, 1.1, "-1", "-1.1", nil, nil},
			output: []any{0.0, 1.0, 1.1, -1.0, -1.1, nil, nil},
		},
		{
			input:  []any{},
			output: []any{},
		},
	}
}

// func objectTests(rnd randomTestValues) []inputOutput {
// 	return []inputOutput{
// 		{
// 			input:  map[string]any{"boolean": true, "integer": 0, "string": "test_string", "float": 0.0},
// 			output: map[string]any{"boolean": true, "integer": 0, "string": "test_string", "float": 0.0},
// 		},
// 		{
// 			input:  map[string]any{"boolean": false, "integer": 1, "string": "3", "float": 1.0},
// 			output: map[string]any{"boolean": false, "integer": 1, "string": "3", "float": 1.0},
// 		},
// 		{
// 			input:  map[string]any{"boolean": true, "integer": 1.0, "string": "3.0", "float": 1.1},
// 			output: map[string]any{"boolean": true, "integer": 1, "string": "3.0", "float": 1.1},
// 		},
// 		{
// 			input:  map[string]any{"boolean": true, "integer": rnd.int, "string": rnd.string, "float": rnd.float63},
// 			output: map[string]any{"boolean": true, "integer": rnd.int, "string": rnd.string, "float": rnd.float63},
// 		},
// 	}
// }

// func dateTypeTests() []inputOutput {
// 	return []inputOutput{
// 		{input: "2021-01-01", output: "2021-01-01"},
// 	}
// }

// func timestampTypeTests() []inputOutput {
// 	return []inputOutput{
// 		{input: "00:00:00", output: "00:00:00"},
// 		{input: "00:00:00.000", output: "00:00:00"},
// 	}
// }

// func dateTimeTypeTests() []inputOutput {
// 	return []inputOutput{
// 		{input: "2021-01-01T00:00:00Z", output: "2021-01-01 00:00:00"},
// 		{input: "2021-01-01T00:00:00.000Z", output: "2021-01-01 00:00:00"},
// 	}
// }

func validateTests() map[string][]inputOutput {
	rnd := generateRandomTestValues()
	enumString, _ := EnumOf(StringType(), stringEnumValues()...)
	floatEnum, _ := EnumOf(FloatType(), floatEnumValues()...)
	// objectType, _ := ObjectType(objectFields())

	return map[string][]inputOutput{
		nonExistingTestField:            nonExistingFieldTests(rnd),
		BooleanType().String():          booleanTypeTests(),
		IntegerType().String():          integerTypeTests(rnd),
		StringType().String():           stringTypeTests(rnd),
		FloatType().String():            floatTypeTests(rnd),
		enumString.String():             stringEnumTests(),
		ArrayOf(BooleanType()).String(): arrayOfBooleanTypeTests(),
		ArrayOf(IntegerType()).String(): arrayOfIntegerTypeTests(rnd),
		ArrayOf(StringType()).String():  arrayOfStringTypeTests(rnd),
		ArrayOf(FloatType()).String():   arrayOfFloatTypeTests(rnd),
		ArrayOf(floatEnum).String():     arrayOfFloatEnumTests(),
		// objectType.String():             objectTests(rnd),
		// DateType().String():             dateTypeTests(),
		// TimestampType().String():        timestampTypeTests(),
		// DateTimeType().String():         dateTimeTypeTests(),
	}
}

func randomString() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	s := make([]rune, rand.Intn(maxStringLength))
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}

func randomBoolean() bool {
	return rand.Intn(maxBoolean) == 0
}

func compareValues(value1, value2 any) bool {
	if value1 == nil && value2 == nil {
		return true
	}

	if value1 == nil || value2 == nil {
		return false
	}

	value1 = transformVal(value1)
	value2 = transformVal(value2)

	val1 := reflect.ValueOf(value1)
	val2 := reflect.ValueOf(value2)

	if val1.Kind() != val2.Kind() {
		return false
	}

	switch val1.Kind() {
	case reflect.Slice:
		return compareSlices(val1, val2)
	case reflect.Map:
		return compareMaps(val1, val2)
	default:
		return value1 == value2
	}
}

func transformVal(val any) any {
	switch val := val.(type) {
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case float32:
		return float64(val)
	default:
		return val
	}
}

func compareSlices(slice1, slice2 reflect.Value) bool {
	if slice1.Len() != slice2.Len() {
		return false
	}

	for i := range slice1.Len() {
		if !compareValues(slice1.Index(i).Interface(), slice2.Index(i).Interface()) {
			return false
		}
	}

	return true
}

func compareMaps(map1, map2 reflect.Value) bool {
	if map1.Len() != map2.Len() {
		return false
	}

	for _, key := range map1.MapKeys() {
		if !compareValues(map1.MapIndex(key).Interface(), map2.MapIndex(key).Interface()) {
			return false
		}
	}

	return true
}
