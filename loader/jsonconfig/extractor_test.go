package jsonconfig

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
	"github.com/xeipuuv/gojsonschema"
)

type ExtractorSuite struct{}

func (s *ExtractorSuite) TestUnmarshalJQValueExtractor(t sweet.T) {
	extractor, err := unmarshalValueExtractor([]byte(`{
		"type": "jq",
		"expr": ".foo.bar"
	}`))

	Expect(err).To(BeNil())
	Expect(extractor.JQ).To(Equal(".foo.bar"))
}

func (s *ExtractorSuite) TestUnmarshalJQValueExtractorList(t sweet.T) {
	extractor, err := unmarshalValueExtractor([]byte(`{
		"type": "jq",
		"expr": ".[].foo",
		"list": true
	}`))

	Expect(err).To(BeNil())
	Expect(extractor.JQ).To(Equal(".[].foo"))
	Expect(extractor.IsList).To(BeTrue())
}

func (s *ExtractorSuite) TestUnmarshalJQValueExtractorWithAssertion(t sweet.T) {
	extractor, err := unmarshalValueExtractor([]byte(`{
		"type": "jq",
		"expr": ".[].foo",
		"assert": {
			"type": "regex",
			"pattern": "^.{3}$"
		}
	}`))

	Expect(err).To(BeNil())
	Expect(extractor.JQ).To(Equal(".[].foo"))
	Expect(extractor.Assert).NotTo(BeNil())
	Expect(extractor.Assert.Pattern.MatchString("foo")).To(BeTrue())
	Expect(extractor.Assert.Pattern.MatchString("bonk")).To(BeFalse())
}

func (s *ExtractorSuite) TestUnmarshalRegexValueExtractor(t sweet.T) {
	extractor, err := unmarshalValueExtractor([]byte(`{
		"type": "regex",
		"pattern": "\\d{4}"
	}`))

	Expect(err).To(BeNil())
	Expect(extractor.Pattern.MatchString("1234")).To(BeTrue())
	Expect(extractor.Pattern.MatchString("abcd")).To(BeFalse())
}

func (s *ExtractorSuite) TestUnmarshalRegexValueExtractorWithHeader(t sweet.T) {
	extractor, err := unmarshalValueExtractor([]byte(`{
		"type": "regex",
		"pattern": "text/(html|xml)",
		"header": "Content-Type"
	}`))

	Expect(err).To(BeNil())
	Expect(extractor.Header).To(Equal("Content-Type"))
}

func (s *ExtractorSuite) TestUnmarshalRegexValueExtractorIllegalRegex(t sweet.T) {
	_, err := unmarshalValueExtractor([]byte(`{
		"type": "regex",
		"pattern": "[abc"
	}`))

	Expect(err).To(MatchError("illegal extractor regex"))
}

func (s *ExtractorSuite) TestUnmarshalUnknownValueExtractor(t sweet.T) {
	_, err := unmarshalValueExtractor([]byte(`{"type": "unknown"}`))
	Expect(err).To(MatchError("unknown extractor type 'unknown'"))
}

func (s *ExtractorSuite) TestUnmarshalRegexValueAssertion(t sweet.T) {
	assertion, err := unmarshalValueAssertion([]byte(`{
		"type": "regex",
		"pattern": "\\d{4}"
	}`))

	Expect(err).To(BeNil())
	Expect(assertion.Pattern.MatchString("1234")).To(BeTrue())
	Expect(assertion.Pattern.MatchString("abcd")).To(BeFalse())
}

func (s *ExtractorSuite) TestUnmarshalRegexValueAssertionIllegalRegex(t sweet.T) {
	_, err := unmarshalValueAssertion([]byte(`{
		"type": "regex",
		"pattern": "[abc"
	}`))

	Expect(err).To(MatchError("illegal assertion regex"))
}

func (s *ExtractorSuite) TestUnmarshalJSONSchemaValueAssertion(t sweet.T) {
	assertion, err := unmarshalValueAssertion([]byte(`{
		"type": "jsonschema",
		"schema": "{\"type\": \"array\"}"
	}`))

	Expect(err).To(BeNil())
	result, err := assertion.Schema.Validate(gojsonschema.NewStringLoader(`{}`))
	Expect(err).To(BeNil())
	Expect(result.Valid()).To(BeFalse())
	Expect(result.Errors()).To(HaveLen(1))
	Expect(result.Errors()[0].Description()).To(Equal("Invalid type. Expected: array, given: object"))
}

func (s *ExtractorSuite) TestUnmarshalJSONSchemaValueAssertionFromYAML(t sweet.T) {
	assertion, err := unmarshalValueAssertion([]byte(`{
		"type": "jsonschema",
		"schema": "type: array"
	}`))

	Expect(err).To(BeNil())
	result, err := assertion.Schema.Validate(gojsonschema.NewStringLoader(`{}`))
	Expect(err).To(BeNil())
	Expect(result.Valid()).To(BeFalse())
	Expect(result.Errors()).To(HaveLen(1))
	Expect(result.Errors()[0].Description()).To(Equal("Invalid type. Expected: array, given: object"))
}

func (s *ExtractorSuite) TestUnmarshalUnknownValueAssertion(t sweet.T) {
	_, err := unmarshalValueAssertion([]byte(`{"type": "unknown"}`))
	Expect(err).To(MatchError("unknown assertion type 'unknown'"))
}
