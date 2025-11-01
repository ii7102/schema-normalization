package rules

import "fmt"

type Field string

func (f Field) String() string {
	return string(f)
}

type BaseType int

const (
	Boolean BaseType = iota
	Integer
	String
	Float
)

var fieldTypeNames = map[BaseType]string{
	Boolean: "boolean",
	Integer: "integer",
	String:  "string",
	Float:   "float",
}

func (bt BaseType) String() string {
	return fieldTypeNames[bt]
}

type FieldType struct {
	BaseType
	IsArray    bool
	EnumValues []any
}

func (ft *FieldType) SetIsArray(isArray bool) {
	ft.IsArray = isArray
}

func (ft *FieldType) SetEnumValues(enumValues ...any) {
	ft.EnumValues = enumValues
}

func (ft *FieldType) AddEnumValue(enumValue any) {
	ft.EnumValues = append(ft.EnumValues, enumValue)
}

func (ft FieldType) String() string {
	str := ft.BaseType.String()
	if ft.EnumValues != nil {
		str = fmt.Sprintf("enum<%s>", str)
	}
	if ft.IsArray {
		str = fmt.Sprintf("array<%s>", str)
	}
	return str
}

func (ft FieldType) Field() Field {
	return Field(ft.String())
}

func BooleanType() FieldType {
	return FieldType{BaseType: Boolean}
}

func IntegerType() FieldType {
	return FieldType{BaseType: Integer}
}

func StringType() FieldType {
	return FieldType{BaseType: String}
}

func FloatType() FieldType {
	return FieldType{BaseType: Float}
}

func ArrayOf(fieldType FieldType) FieldType {
	fieldType.SetIsArray(true)
	return fieldType
}

func EnumOf(fieldType FieldType, enumValues ...any) FieldType {
	fieldType.SetEnumValues(enumValues...)
	return fieldType
}
