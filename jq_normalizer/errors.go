package jqnormalizer

import "errors"

var (
	errJqQueryCompilationNotCached = errors.New("jq query compilation not cached")
	errNoResultFromJqQuery         = errors.New("no result from JQ query")

	errUnexpectedResultType = errors.New("unexpected result type")
)
