package schema

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ii7102/schema-normalization/errors"
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
	Timestamp
	Date
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

// LayoutFormat returns the layout format for the given base type.
func (bt BaseType) LayoutFormat() (string, error) {
	switch bt {
	case Timestamp:
		return time.TimeOnly, nil
	case Date:
		return time.DateOnly, nil
	case DateTime:
		return time.DateTime, nil
	default:
		return "", errors.WrappedError(errBaseTypeWithoutLayout, "base type %s does not have a layout", bt)
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
	layout   *string
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

// Layout returns the layout for the given fieldType.
func (ft FieldType) Layout() *string {
	return ft.layout
}

// SetLayout validates the given layout and sets it to the given fieldType.
func (ft *FieldType) SetLayout(layout string) error {
	if err := validateDateTime(layout, ft.baseType); err != nil {
		return err
	}

	ft.layout = &layout

	return nil
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

// SetObjectFields validates the given object fields and sets them to the given fieldType.
func (ft *FieldType) SetObjectFields(objectFields map[Field]FieldType) error {
	if ft.baseType != Object {
		return errors.WrappedError(errFieldTypeIsNotAnObject, "field type %s", ft.baseType)
	}

	if err := validateObjectFields(objectFields); err != nil {
		return err
	}

	ft.object = &object{fields: objectFields}

	return nil
}

// AddObjectField validates the given field and fieldType and adds them to the given fieldType.
func (ft *FieldType) AddObjectField(field Field, fieldType FieldType) error {
	if ft.object == nil || ft.object.fields == nil {
		return errors.WrappedError(errObjectFieldsNotSet, "field type %s", ft.baseType)
	}

	if ft.baseType != Object {
		return errors.WrappedError(errFieldTypeIsNotAnObject, "field type %s", ft.baseType)
	}

	if fieldType.baseType == Object {
		return errors.WrappedError(errNestedObjectsAreNotSupported, "object field '%s' cannot be of type object ", field)
	}

	ft.object.fields[field] = fieldType

	return nil
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
