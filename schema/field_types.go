package schema

import (
	"time"

	"github.com/ii7102/schema-normalization/errors"
)

// ArrayOf updates the given fieldType to be an array.
func ArrayOf(fieldType FieldType) FieldType {
	fieldType.SetIsArray(true)

	return fieldType
}

// BooleanType returns a new FieldType with the boolean base type set.
func BooleanType() FieldType {
	return FieldType{baseType: Boolean}
}

// IntegerType returns a new FieldType with the integer base type set.
func IntegerType() FieldType {
	return FieldType{baseType: Integer}
}

// StringType returns a new FieldType with the string base type set.
func StringType() FieldType {
	return FieldType{baseType: String}
}

// FloatType returns a new FieldType with the float base type set.
func FloatType() FieldType {
	return FieldType{baseType: Float}
}

// ObjectType validates the given object fields and returns a new FieldType with the object fields set.
func ObjectType(objectFields map[Field]FieldType) (FieldType, error) {
	if err := validateObjectFields(objectFields); err != nil {
		return FieldType{}, errors.WrappedError(err, "invalid object fields: %v", objectFields)
	}

	return FieldType{
		baseType: Object,
		object:   &object{fields: objectFields},
	}, nil
}

// EnumOf sets the given enum values to the given fieldType.
func EnumOf(fieldType FieldType, enumValues ...any) (FieldType, error) {
	if err := fieldType.SetEnumValues(enumValues...); err != nil {
		return FieldType{}, err
	}

	return fieldType, nil
}

// DateType returns a new FieldType with the date base type set.
func DateType() FieldType {
	layout := time.DateOnly

	return FieldType{baseType: Date, layout: &layout}
}

// TimestampType validates the given layout and returns a new FieldType with the timestamp base type set.
func TimestampType(layout string) (FieldType, error) {
	if err := validateDateTime(layout, Timestamp); err != nil {
		return FieldType{}, err
	}

	return FieldType{baseType: Timestamp, layout: &layout}, nil
}

// DateTimeType validates the given layout and returns a new FieldType with the dateTime base type set.
func DateTimeType(layout string) (FieldType, error) {
	if err := validateDateTime(layout, DateTime); err != nil {
		return FieldType{}, err
	}

	return FieldType{baseType: DateTime, layout: &layout}, nil
}
