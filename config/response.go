package config

import (
	"regexp"

	"github.com/xeipuuv/gojsonschema"
)

type (
	Response struct {
		Status  *regexp.Regexp
		Extract map[string]*ValueExtractor
	}

	ValueExtractor struct {
		JQ      string
		IsList  bool
		Pattern *regexp.Regexp
		Assert  *ValueAssertion
		Header  string
	}

	ValueAssertion struct {
		Pattern *regexp.Regexp
		Schema  *gojsonschema.Schema
	}
)
