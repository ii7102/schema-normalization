package jqnormalizer

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"diploma/rules"
	"github.com/itchyny/gojq"
)

const (
	commonRule = "\n%[1]s: (.%[1]s | %s )"

	nullRule = "if . == null then . else "

	booleanRule      = "if type == \"string\" then (. == \"true\") else . end"
	integerRule      = nullRule + "tonumber | round end"
	stringRule       = nullRule + "tostring end"
	floatRule        = nullRule + "tonumber * 1.0 end"
	booleanArrayRule = nullRule + "map(if type == \"string\" then (. == \"true\") else . end) end"
	integerArrayRule = nullRule + "map(" + nullRule + "tonumber | round end) end"
	stringArrayRule  = nullRule + "map(" + nullRule + "tostring end) end"
	floatArrayRule   = nullRule + "map(" + nullRule + "tonumber * 1.0 end) end"
)

func jqRules() map[string]string {
	return map[string]string{
		"boolean":        booleanRule,
		"integer":        integerRule,
		"string":         stringRule,
		"float":          floatRule,
		"array<boolean>": booleanArrayRule,
		"array<integer>": integerArrayRule,
		"array<string>":  stringArrayRule,
		"array<float>":   floatArrayRule,
	}
}

func jqFilter(fields map[rules.Field]rules.FieldType) string {
	jqRules := jqRules()

	jqRulesArray := make([]string, 0, len(fields))
	for field, fieldType := range fields {
		jqRule, ok := jqRules[fieldType.String()]
		if ok {
			jqRulesArray = append(jqRulesArray, fmt.Sprintf(commonRule, field, jqRule))
		}
	}

	return fmt.Sprintf("{%s\n}", strings.Join(jqRulesArray, ","))
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
