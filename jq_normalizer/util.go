package jqnormalizer

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

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

	dateRule      = "if test(\"^[0-9]{4}-[0-9]{2}-[0-9]{2}$\") then . else error(\"Invalid date: \" + .) end "
	timestampRule = "try (strptime(%q) | strftime(\"%%H:%%M:%%S\")) catch null "
	dateTimeRule  = "try (strptime(%q) | strftime(\"%%Y-%%m-%%d %%H:%%M:%%S\")) catch null "
)

func timestampStrptimeFormat(layout *string) string {
	if layout == nil {
		return ""
	}

	format, found := map[string]string{
		time.TimeOnly:   "%H:%M:%S",
		time.Kitchen:    "%I:%M%p",
		time.Stamp:      "%b %e %H:%M:%S",
		time.StampMilli: "%b %e %H:%M:%S.%f",
		time.StampMicro: "%b %e %H:%M:%S.%f",
		time.StampNano:  "%b %e %H:%M:%S.%f",
	}[*layout]

	if !found {
		return ""
	}

	return format
}

func dateTimeStrptimeFormat(layout *string) string {
	if layout == nil {
		return ""
	}

	format, found := map[string]string{
		time.RFC3339:     "%Y-%m-%dT%H:%M:%S%z",
		time.DateTime:    "%Y-%m-%d %H:%M:%S",
		time.RFC3339Nano: "%Y-%m-%dT%H:%M:%S.%f%z",
		time.RFC1123:     "%a, %d %b %Y %H:%M:%S %Z",
		time.RFC1123Z:    "%a, %d %b %Y %H:%M:%S %z",
		time.RFC850:      "%A, %d-%b-%y %H:%M:%S %Z",
		time.RFC822:      "%d %b %y %H:%M %Z",
		time.RFC822Z:     "%d %b %y %H:%M %z",
		time.ANSIC:       "%a %b %e %H:%M:%S %Y",
		time.UnixDate:    "%a %b %e %H:%M:%S %Z %Y",
		time.RubyDate:    "%a %b %d %H:%M:%S %z %Y",
	}[*layout]

	if !found {
		return ""
	}

	return format
}

func baseTypeJqRule(fieldType rules.FieldType) string {
	return map[rules.BaseType]string{
		rules.Boolean:   booleanRule,
		rules.Integer:   integerRule,
		rules.String:    stringRule,
		rules.Float:     floatRule,
		rules.Date:      dateRule,
		rules.Timestamp: fmt.Sprintf(timestampRule, timestampStrptimeFormat(fieldType.Layout())),
		rules.DateTime:  fmt.Sprintf(dateTimeRule, dateTimeStrptimeFormat(fieldType.Layout())),
		rules.Object:    jqFilter(fieldType.ObjectFields()),
	}[fieldType.BaseType()]
}

func generateJqRule(fieldType rules.FieldType) string {
	var (
		jqRule  = baseTypeJqRule(fieldType)
		isArray = fieldType.IsArray()
		enums   = fieldType.EnumValues()
	)

	if jqRule == "" {
		return ""
	}

	switch {
	case isArray && len(enums) > 0:
		jqRule = fmt.Sprintf(arrayEnumRule, enums, jqRule)
	case len(enums) > 0:
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
