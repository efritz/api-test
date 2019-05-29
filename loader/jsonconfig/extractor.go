package jsonconfig

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/efritz/api-test/config"
	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

type (
	TypeEnvelope struct {
		Type string `json:"type"`
	}

	JQValueExtractor struct {
		Expr   string          `json:"expr"`
		IsList bool            `json:"list"`
		Assert json.RawMessage `json:"assert"`
		Header string          `json:"header"`
	}

	RegexValueExtractor struct {
		Pattern string          `json:"pattern"`
		Assert  json.RawMessage `json:"assert"`
		Header  string          `json:"header"`
	}

	RegexValueAssertion struct {
		Pattern string `json:"pattern"`
	}

	JSONSchemaValueAssertion struct {
		Schema string `json:"schema"`
	}
)

var (
	unmarshalExtractFuncs = map[string]func(json.RawMessage) (*config.ValueExtractor, error){
		"jq":    unmarshalJQValueExtractor,
		"regex": unmarshalRegexValueExtractor,
	}

	unmarshalAssertionFuncs = map[string]func(json.RawMessage) (*config.ValueAssertion, error){
		"regex":      unmarshalRegexValueAssertion,
		"jsonschema": unmarshalJSONSchemaValueAssertion,
	}
)

//
// Extractors

func unmarshalValueExtractor(payload json.RawMessage) (*config.ValueExtractor, error) {
	envelope := &TypeEnvelope{}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return nil, err
	}

	fn, ok := unmarshalExtractFuncs[envelope.Type]
	if !ok {
		return nil, fmt.Errorf("unknown extractor type '%s'", envelope.Type)
	}

	return fn(payload)
}

func unmarshalJQValueExtractor(payload json.RawMessage) (*config.ValueExtractor, error) {
	properties := &JQValueExtractor{}
	if err := json.Unmarshal(payload, &properties); err != nil {
		return nil, err
	}

	assertion, err := unmarshalValueAssertion(properties.Assert)
	if err != nil {
		return nil, err
	}

	return &config.ValueExtractor{
		JQ:     properties.Expr,
		IsList: properties.IsList,
		Assert: assertion,
		Header: properties.Header,
	}, nil
}

func unmarshalRegexValueExtractor(payload json.RawMessage) (*config.ValueExtractor, error) {
	properties := &RegexValueExtractor{}
	if err := json.Unmarshal(payload, &properties); err != nil {
		return nil, err
	}

	pattern, err := regexp.Compile(properties.Pattern)
	if err != nil {
		return nil, fmt.Errorf("illegal extractor regex")
	}

	assertion, err := unmarshalValueAssertion(properties.Assert)
	if err != nil {
		return nil, err
	}

	return &config.ValueExtractor{
		Pattern: pattern,
		Assert:  assertion,
		Header:  properties.Header,
	}, nil
}

//
// Assertions

func unmarshalValueAssertion(payload json.RawMessage) (*config.ValueAssertion, error) {
	if payload == nil {
		return nil, nil
	}

	envelope := &TypeEnvelope{}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return nil, err
	}

	fn, ok := unmarshalAssertionFuncs[envelope.Type]
	if !ok {
		return nil, fmt.Errorf("unknown assertion type '%s'", envelope.Type)
	}

	return fn(payload)
}

func unmarshalRegexValueAssertion(payload json.RawMessage) (*config.ValueAssertion, error) {
	properties := &RegexValueAssertion{}
	if err := json.Unmarshal(payload, &properties); err != nil {
		return nil, err
	}

	pattern, err := regexp.Compile(properties.Pattern)
	if err != nil {
		return nil, fmt.Errorf("illegal assertion regex")
	}

	return &config.ValueAssertion{
		Pattern: pattern,
	}, nil
}

func unmarshalJSONSchemaValueAssertion(payload json.RawMessage) (*config.ValueAssertion, error) {
	properties := &JSONSchemaValueAssertion{}
	if err := json.Unmarshal(payload, &properties); err != nil {
		return nil, err
	}

	rawSchema, err := yaml.YAMLToJSON([]byte(properties.Schema))
	if err != nil {
		return nil, fmt.Errorf("failed to convert schema to JSON (%s)", err.Error())
	}

	schema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(rawSchema))
	if err != nil {
		return nil, fmt.Errorf("Failed to load schema (%s)", err.Error())
	}

	return &config.ValueAssertion{
		Schema: schema,
	}, nil
}
