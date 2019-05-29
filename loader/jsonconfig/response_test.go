package jsonconfig

import (
	"encoding/json"
	"fmt"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ResponseSuite struct{}

func (s *ResponseSuite) TestTranslate(t sweet.T) {
	response := &Response{
		Status: json.RawMessage(`"2.."`),
		Extract: map[string]json.RawMessage{
			"foo": json.RawMessage([]byte(`{"type": "jq", "expr": ".foo"}`)),
			"bar": json.RawMessage([]byte(`{"type": "jq", "expr": ".bar"}`)),
			"baz": json.RawMessage([]byte(`{"type": "jq", "expr": ".baz"}`)),
		},
	}

	translated, err := response.Translate()
	Expect(err).To(BeNil())
	Expect(translated).NotTo(BeNil())
	Expect(translated.Status.MatchString("200")).To(BeTrue())
	Expect(translated.Status.MatchString("300")).To(BeFalse())

	for _, name := range []string{"foo", "bar", "baz"} {
		Expect(translated.Extract).To(HaveKey(name))
		Expect(translated.Extract[name].JQ).To(Equal(fmt.Sprintf(".%s", name)))
	}
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

// func (s *ResponseSuite) TestTranslateStringLists(t sweet.T) {
// 	response := &Response{
// 		Headers: map[string]json.RawMessage{
// 			"X-Test": json.RawMessage(`["foo", "bar"]`),
// 		},
// 	}

// 	translated, err := response.Translate()
// 	Expect(err).To(BeNil())
// 	Expect(translated).NotTo(BeNil())
// 	Expect(translated.Headers).To(HaveKey("X-Test"))
// 	Expect(translated.Headers["X-Test"]).To(HaveLen(2))
// 	Expect(translated.Headers["X-Test"][0].MatchString("foo")).To(BeTrue())
// 	Expect(translated.Headers["X-Test"][0].MatchString("baz")).To(BeFalse())
// 	Expect(translated.Headers["X-Test"][1].MatchString("bar")).To(BeTrue())
// 	Expect(translated.Headers["X-Test"][1].MatchString("baz")).To(BeFalse())
// }

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
