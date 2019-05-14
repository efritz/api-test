package jsonconfig

import (
	"encoding/json"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type RequestSuite struct{}

func (s *RequestSuite) TestTranslate(t sweet.T) {
	request := &Request{
		URI:    "/users",
		Method: "post",
		Auth: &BasicAuth{
			Username: "admin",
			Password: "secret",
		},
		Headers: map[string]json.RawMessage{
			"X-Custom1": []byte(`"foobar"`),
			"X-Custom2": []byte(`"barbaz"`),
		},
		Body: "payload",
	}

	translated, err := request.Translate(nil)
	Expect(err).To(BeNil())
	Expect(testExec(translated.URL)).To(Equal("/users"))
	Expect(translated.Method).To(Equal("post"))
	Expect(translated.Auth).NotTo(BeNil())
	Expect(testExec(translated.Auth.Username)).To(Equal("admin"))
	Expect(testExec(translated.Auth.Password)).To(Equal("secret"))
	Expect(translated.Headers).To(HaveKey("X-Custom1"))
	Expect(translated.Headers).To(HaveKey("X-Custom2"))
	Expect(translated.Headers["X-Custom1"]).To(HaveLen(1))
	Expect(translated.Headers["X-Custom2"]).To(HaveLen(1))
	Expect(testExec(translated.Headers["X-Custom1"][0])).To(Equal("foobar"))
	Expect(testExec(translated.Headers["X-Custom2"][0])).To(Equal("barbaz"))
	Expect(testExec(translated.Body)).To(Equal("payload"))
}

func (s *RequestSuite) TestTranslateWithJSONBody(t sweet.T) {
	request := &Request{
		JSONBody: []byte(`{"x": 1, "y": 2, "z": 3}`),
	}

	translated, err := request.Translate(nil)
	Expect(err).To(BeNil())
	Expect(testExec(translated.Body)).To(Equal(`{"x": 1, "y": 2, "z": 3}`))
}

func (s *RequestSuite) TestTranslateNoExplicitMethod(t sweet.T) {
	request := &Request{}
	translated, err := request.Translate(nil)
	Expect(err).To(BeNil())
	Expect(translated.Method).To(Equal("get"))
}

func (s *RequestSuite) TestTranslateNoAuth(t sweet.T) {
	request := &Request{}
	translated, err := request.Translate(nil)
	Expect(err).To(BeNil())
	Expect(translated.Auth).To(BeNil())
}

func (s *RequestSuite) TestTranslateGlobalRequestURL(t sweet.T) {
	request := &Request{URI: "/users"}
	translated, err := request.Translate(&GlobalRequest{
		BaseURL: "http://test.io",
	})

	Expect(err).To(BeNil())
	Expect(testExec(translated.URL)).To(Equal("http://test.io/users"))
}

func (s *RequestSuite) TestTranslateGlobalRequestURLAbsoluteURI(t sweet.T) {
	request := &Request{URI: "http://test.io/users"}
	translated, err := request.Translate(&GlobalRequest{
		BaseURL: "http://wrong.io",
	})

	Expect(err).To(BeNil())
	Expect(testExec(translated.URL)).To(Equal("http://test.io/users"))
}

func (s *RequestSuite) TestTranslateGlobalRequestAuth(t sweet.T) {
	request := &Request{}
	translated, err := request.Translate(&GlobalRequest{
		Auth: &BasicAuth{
			Username: "admin",
			Password: "secret",
		},
	})

	Expect(err).To(BeNil())
	Expect(translated.Auth).NotTo(BeNil())
	Expect(testExec(translated.Auth.Username)).To(Equal("admin"))
	Expect(testExec(translated.Auth.Password)).To(Equal("secret"))
}

func (s *RequestSuite) TestTranslateGlobalRequestAuthOverride(t sweet.T) {
	request := &Request{
		Auth: &BasicAuth{
			Username: "adminer",
			Password: "secreter",
		},
	}

	translated, err := request.Translate(&GlobalRequest{
		Auth: &BasicAuth{
			Username: "admin",
			Password: "secret",
		},
	})

	Expect(err).To(BeNil())
	Expect(translated.Auth).NotTo(BeNil())
	Expect(testExec(translated.Auth.Username)).To(Equal("adminer"))
	Expect(testExec(translated.Auth.Password)).To(Equal("secreter"))
}

func (s *RequestSuite) TestTranslateGlobalRequestHeaders(t sweet.T) {
	request := &Request{
		Headers: map[string]json.RawMessage{
			"X-Custom2": []byte(`"baz"`),
		},
	}

	translated, err := request.Translate(&GlobalRequest{
		Headers: map[string]json.RawMessage{
			"X-Custom1": []byte(`["foo", "bar"]`),
			"X-Custom2": []byte(`"bonk"`),
		},
	})

	Expect(err).To(BeNil())
	Expect(translated.Headers).To(HaveKey("X-Custom1"))
	Expect(translated.Headers).To(HaveKey("X-Custom2"))
	Expect(translated.Headers["X-Custom1"]).To(HaveLen(2))
	Expect(translated.Headers["X-Custom2"]).To(HaveLen(1))
	Expect(testExec(translated.Headers["X-Custom1"][0])).To(Equal("foo"))
	Expect(testExec(translated.Headers["X-Custom1"][1])).To(Equal("bar"))
	Expect(testExec(translated.Headers["X-Custom2"][0])).To(Equal("baz"))
}

func (s *RequestSuite) TestTranslateStringLists(t sweet.T) {
	request := &Request{
		Headers: map[string]json.RawMessage{
			"X-Custom1": []byte(`["foo", "bar"]`),
			"X-Custom2": []byte(`["bar", "baz"]`),
		},
	}

	translated, err := request.Translate(&GlobalRequest{
		Headers: map[string]json.RawMessage{
			"X-Custom3": []byte(`["baz", "bonk"]`),
		},
	})

	Expect(err).To(BeNil())
	Expect(translated.Headers).To(HaveKey("X-Custom1"))
	Expect(translated.Headers).To(HaveKey("X-Custom2"))
	Expect(translated.Headers).To(HaveKey("X-Custom3"))
	Expect(translated.Headers["X-Custom1"]).To(HaveLen(2))
	Expect(translated.Headers["X-Custom2"]).To(HaveLen(2))
	Expect(translated.Headers["X-Custom3"]).To(HaveLen(2))
	Expect(testExec(translated.Headers["X-Custom1"][0])).To(Equal("foo"))
	Expect(testExec(translated.Headers["X-Custom1"][1])).To(Equal("bar"))
	Expect(testExec(translated.Headers["X-Custom2"][0])).To(Equal("bar"))
	Expect(testExec(translated.Headers["X-Custom2"][1])).To(Equal("baz"))
	Expect(testExec(translated.Headers["X-Custom3"][0])).To(Equal("baz"))
	Expect(testExec(translated.Headers["X-Custom3"][1])).To(Equal("bonk"))
}

func (s *RequestSuite) TestTranslateMutuallyExclusiveBodies(t sweet.T) {
	request := &Request{
		Body:     "body",
		JSONBody: json.RawMessage(`["another", "body"]`),
	}

	_, err := request.Translate(nil)
	Expect(err).To(MatchError("multiple bodies supplied"))
}

func (s *RequestSuite) TestTranslateInvalidURITemplate(t sweet.T) {
	request := &Request{
		URI: "{{",
	}

	_, err := request.Translate(nil)
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("illegal uri template"))
}

func (s *RequestSuite) TestTranslateInvalidHeaderTemplate(t sweet.T) {
	request := &Request{
		Headers: map[string]json.RawMessage{
			"X-Request-ID": json.RawMessage(`"{{"`),
		},
	}

	_, err := request.Translate(nil)
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("illegal header template"))
}

func (s *RequestSuite) TestTranslateInvalidBodyTemplate(t sweet.T) {
	request := &Request{
		Body: "{{",
	}

	_, err := request.Translate(nil)
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("illegal body template"))
}

func (s *RequestSuite) TestTranslateInvalidJSONBodyTemplate(t sweet.T) {
	request := &Request{
		JSONBody: json.RawMessage(`{"foo": "{{"}`),
	}

	_, err := request.Translate(nil)
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("illegal json body template"))
}

func (s *RequestSuite) TestCompileUUID(t sweet.T) {
	template, err := compile("{{uuid}}")
	Expect(err).To(BeNil())

	uuid1 := testExec(template)
	uuid2 := testExec(template)
	uuid3 := testExec(template)
	Expect(uuid1).NotTo(Equal(uuid2))
	Expect(uuid1).NotTo(Equal(uuid3))
	Expect(uuid2).NotTo(Equal(uuid3))
	Expect(uuid1).To(MatchRegexp("[a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12}"))
}

func (s *RequestSuite) TestCompileFile(t sweet.T) {
	template, err := compile(`{{file "request_test.go"}}`)
	Expect(err).To(BeNil())
	Expect(testExec(template)).To(HavePrefix("package jsonconfig\n"))
}
