package rules

import (
	"fmt"
	"strconv"
	"strings"
)

// Field is a string representation of a field name.
type Field string

func (f Field) String() string {
	return string(f)
}

// BaseType is an enum of the base types supported by the normalizer.
type BaseType int

// Supported enum values for the BaseType.
const (
	Boolean BaseType = iota
	Integer
	String
	Float
	Date
	Timestamp
	DateTime
	Object
)

func (bt BaseType) String() string {
	switch bt {
	case Boolean:
		return "boolean"
	case Integer:
		return "integer"
	case String:
		return "string"
	case Float:
		return "float"
	case Date:
		return "date"
	case Timestamp:
		return "timestamp"
	case DateTime:
		return "dateTime"
	case Object:
		return "object"
	default:
		return ""
	}
}

type enumValues []any

func (ev enumValues) String() string {
	if len(ev) == 0 {
		return ""
	}

	enumValues := make([]string, len(ev))
	for i, enumValue := range ev {
		if s, ok := enumValue.(string); ok {
			enumValues[i] = strconv.Quote(s)
		} else {
			enumValues[i] = fmt.Sprintf("%v", enumValue)
		}
	}

	return fmt.Sprintf("[%s]", strings.Join(enumValues, ", "))
}

type enum struct {
	values enumValues
}

type object struct {
	fields map[Field]FieldType
}

// FieldType represents a field type with its base type and additional properties.
type FieldType struct {
	baseType BaseType
	isArray  bool
	enum     *enum
	object   *object
}

// BaseType returns the base type for the given fieldType.
func (ft FieldType) BaseType() BaseType {
	return ft.baseType
}

// IsArray return true if the given fieldType is an array and false otherwise.
func (ft FieldType) IsArray() bool {
	return ft.isArray
}

// SetIsArray sets the given isArray to the given fieldType.
func (ft *FieldType) SetIsArray(isArray bool) {
	ft.isArray = isArray
}

// EnumValues returns the enum values for the given fieldType.
func (ft FieldType) EnumValues() enumValues {
	if ft.enum == nil {
		return nil
	}

	return ft.enum.values
}

// SetEnumValues validates the given enum values and sets them to the given fieldType.
func (ft *FieldType) SetEnumValues(enumValues ...any) error {
	if err := validateEnumValues(ft.baseType, enumValues...); err != nil {
		return err
	}

	ft.enum = &enum{values: enumValues}

	return nil
}

// AddEnumValue adds the given enum value to the given fieldType.
// It returns an error if the enum values are not set.
func (ft *FieldType) AddEnumValue(enumValue any) error {
	if ft.enum == nil {
		return errEnumValuesNotSet
	}

	if err := validateEnumValues(ft.baseType, enumValue); err != nil {
		return err
	}

	ft.enum.values = append(ft.enum.values, enumValue)

	return nil
}

// ObjectFields returns the object fields for the given fieldType.
func (ft FieldType) ObjectFields() map[Field]FieldType {
	if ft.object == nil {
		return nil
	}

	return ft.object.fields
}

func (ft FieldType) String() string {
	str := ft.baseType.String()
	if ft.enum != nil {
		str = "enum_" + str
	}

	if ft.isArray {
		str = "array_" + str
	}

	return str
}

// Field returns the field name for the given fieldType.
func (ft FieldType) Field() Field {
	return Field(ft.String())
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

// DateType returns a new FieldType with the date base type set.
func DateType() FieldType {
	return FieldType{baseType: Date}
}

// TimestampType returns a new FieldType with the timestamp base type set.
func TimestampType() FieldType {
	return FieldType{baseType: Timestamp}
}

// DateTimeType returns a new FieldType with the dateTime base type set.
func DateTimeType() FieldType {
	return FieldType{baseType: DateTime}
}

// ObjectType validates the given object fields and returns a new FieldType with the object fields set.
func ObjectType(objectFields map[Field]FieldType) (FieldType, error) {
	if err := validateObjectFields(objectFields); err != nil {
		return FieldType{}, WrappedError(err, "invalid object fields: %v", objectFields)
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

// ArrayOf updates the given fieldType to be an array.
func ArrayOf(fieldType FieldType) FieldType {
	fieldType.SetIsArray(true)

	return fieldType
}
