package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type TestSuite struct{}

func (s *TestSuite) TestTranslate(t sweet.T) {
	enabled := true

	test := &Test{
		Name:    "foobar",
		Enabled: &enabled,
		Request: &Request{
			URI: "/users",
		},
		Response: &Response{
			Status: json.RawMessage(`"2.."`),
		},
	}

	translated, err := test.Translate(&GlobalRequest{
		BaseURL: "http://test.io",
	})

	Expect(err).To(BeNil())
	Expect(translated.Name).To(Equal("foobar"))
	Expect(translated.Enabled).To(BeTrue())
	Expect(testExec(translated.Request.URL)).To(Equal("http://test.io/users"))
	Expect(translated.Response.Status.MatchString("204")).To(BeTrue())
}

func (s *TestSuite) TestTranslateDefaultName(t sweet.T) {
	test := &Test{Request: &Request{URI: "/users", Method: "post"}}
	translated, err := test.Translate(nil)
	Expect(err).To(BeNil())
	Expect(translated.Name).To(Equal("POST /users"))
}

func (s *TestSuite) TestTranslateDefaultEnabled(t sweet.T) {
	test := &Test{Enabled: nil}
	translated, err := test.Translate(nil)
	Expect(err).To(BeNil())
	Expect(translated.Enabled).To(BeTrue())
}

func (s *TestSuite) TestTranslateNoResponse(t sweet.T) {
	test := &Test{Enabled: nil}
	translated, err := test.Translate(nil)
	Expect(err).To(BeNil())
	Expect(translated.Response.Status.MatchString("204")).To(BeTrue())
}
