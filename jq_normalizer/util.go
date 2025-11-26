package jqnormalizer

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ii7102/schema-normalization/rules"
	"github.com/itchyny/gojq"
)

const (
	commonRule = "\n%[1]s: .%[1]s | %s"

	nullRule    = "if . == null then null else %s"
	endRule     = "end "
	nullEndRule = nullRule + endRule

	booleanRule = "if . == \"true\" then true elif . == \"false\" then false else . " + endRule
	integerRule = "tonumber | trunc "
	stringRule  = "tostring "
	floatRule   = "tonumber "

	enumErrorRule = "error(\"Invalid enum value: \" + $val | tostring) "
	enumCheck     = "as $val | if $enums | index($val) then $val else " + enumErrorRule + endRule

	enumRule = "%s as $enums | %s " + enumCheck

	arrayRule = "map( " + nullEndRule + ") "

	arrayEnumRule = "%s as $enums | map( " + nullRule + enumCheck + endRule + ") "
)

func baseTypeJqRule(baseType rules.BaseType) string {
	switch baseType {
	case rules.Boolean:
		return booleanRule
	case rules.Integer:
		return integerRule
	case rules.String:
		return stringRule
	case rules.Float:
		return floatRule
	default:
		return ""
	}
}

func generateJqRule(fieldType rules.FieldType) string {
	var (
		jqRule  = baseTypeJqRule(fieldType.BaseType())
		isArray = fieldType.IsArray()
		enums   = fieldType.EnumValues().String()
	)

	if jqRule == "" {
		return ""
	}

	switch {
	case isArray && enums != "":
		jqRule = fmt.Sprintf(arrayEnumRule, enums, jqRule)
	case enums != "":
		jqRule = fmt.Sprintf(enumRule, enums, jqRule)
	case isArray:
		jqRule = fmt.Sprintf(arrayRule, jqRule)
	default:
	}

	return fmt.Sprintf(nullEndRule, jqRule)
}

func jqFilter(fields map[rules.Field]rules.FieldType) string {
	jqRules := make([]string, 0, len(fields))
	for field, fieldType := range fields {
		if jqRule := generateJqRule(fieldType); jqRule != "" {
			jqRules = append(jqRules, fmt.Sprintf(commonRule, field, jqRule))
		}
	}

	return fmt.Sprintf("{%s\n}", strings.Join(jqRules, ","))
}

func jqBatchFilter(fields map[rules.Field]rules.FieldType) string {
	return fmt.Sprintf("[.[] | %s]", jqFilter(fields))
}

func compileJqCode(jqFilter string) *gojq.Code {
	jqQuery, err := gojq.Parse(jqFilter)
	if err != nil {
		parseErr := &gojq.ParseError{}
		if errors.As(err, &parseErr) {
			log.Printf("JQ parse error at position %d, token %q: %v", parseErr.Offset, parseErr.Token, parseErr.Error())
		} else {
			log.Printf("failed to parse JQ query: %v", err)
		}

		return nil
	}

	compiledCode, err := gojq.Compile(jqQuery)
	if err != nil {
		log.Printf("failed to compile JQ query: %v", err)

		return nil
	}

	return compiledCode
}
