package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ResponseSuite struct{}

func (s *ResponseSuite) TestTranslate(t sweet.T) {
	response := &Response{
		Status: json.RawMessage(`"2.."`),
		Headers: map[string]json.RawMessage{
			"X-Test": json.RawMessage(`"text/.*"`),
		},
		Body: "body",
		ExtractList: map[string]string{
			"ids": ".[].id",
		},
	}

	translated, err := response.Translate()
	Expect(err).To(BeNil())
	Expect(translated).NotTo(BeNil())
	Expect(translated.Status.MatchString("200")).To(BeTrue())
	Expect(translated.Status.MatchString("300")).To(BeFalse())
	Expect(translated.Body.MatchString("body")).To(BeTrue())
	Expect(translated.Body.MatchString("fail")).To(BeFalse())
	Expect(translated.Headers).To(HaveKey("X-Test"))
	Expect(translated.Headers["X-Test"]).To(HaveLen(1))
	Expect(translated.Headers["X-Test"][0].MatchString("text/html")).To(BeTrue())
	Expect(translated.Headers["X-Test"][0].MatchString("data/json")).To(BeFalse())
	Expect(translated.ExtractList).To(HaveKeyWithValue("ids", ".[].id"))
}

func (s *ResponseSuite) TestTranslateNumericStatus(t sweet.T) {
	response := &Response{
		Status: json.RawMessage(`200`),
	}

	translated, err := response.Translate()
	Expect(err).To(BeNil())
	Expect(translated).NotTo(BeNil())
	Expect(translated.Status.MatchString("200")).To(BeTrue())
	Expect(translated.Status.MatchString("204")).To(BeFalse())
}

func (s *ResponseSuite) TestTranslateNoExplicitStatus(t sweet.T) {
	response := &Response{}
	translated, err := response.Translate()
	Expect(err).To(BeNil())
	Expect(translated).NotTo(BeNil())
	Expect(translated.Status.MatchString("200")).To(BeTrue())
	Expect(translated.Status.MatchString("204")).To(BeTrue())
	Expect(translated.Status.MatchString("300")).To(BeFalse())
}

func (s *ResponseSuite) TestTranslateStringLists(t sweet.T) {
	response := &Response{
		Headers: map[string]json.RawMessage{
			"X-Test": json.RawMessage(`["foo", "bar"]`),
		},
	}

	translated, err := response.Translate()
	Expect(err).To(BeNil())
	Expect(translated).NotTo(BeNil())
	Expect(translated.Headers).To(HaveKey("X-Test"))
	Expect(translated.Headers["X-Test"]).To(HaveLen(2))
	Expect(translated.Headers["X-Test"][0].MatchString("foo")).To(BeTrue())
	Expect(translated.Headers["X-Test"][0].MatchString("baz")).To(BeFalse())
	Expect(translated.Headers["X-Test"][1].MatchString("bar")).To(BeTrue())
	Expect(translated.Headers["X-Test"][1].MatchString("baz")).To(BeFalse())
}

func (s *ResponseSuite) TestTranslateInvalidStatusDatatype(t sweet.T) {
	response := &Response{
		Status: json.RawMessage(`[1, 2, 3]`),
	}

	_, err := response.Translate()
	Expect(err).To(MatchError("status value is neither string nor int"))
}

func (s *ResponseSuite) TestTranslateInvalidStatusRegex(t sweet.T) {
	response := &Response{
		Status: json.RawMessage(`"("`),
	}

	_, err := response.Translate()
	Expect(err).To(MatchError("illegal status regex"))
}

func (s *ResponseSuite) TestTranslateInvalidHeaderRegex(t sweet.T) {
	response := &Response{
		Headers: map[string]json.RawMessage{
			"X-Test": json.RawMessage(`"("`),
		},
	}

	_, err := response.Translate()
	Expect(err).To(MatchError("illegal header regex"))
}

func (s *ResponseSuite) TestTranslateInvalidBodyRegex(t sweet.T) {
	response := &Response{
		Body: "(",
	}

	_, err := response.Translate()
	Expect(err).To(MatchError("illegal body regex"))
}
