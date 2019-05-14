package runner

import (
	"net/http"
	tmpl "text/template"

	"github.com/aphistic/sweet"
	"github.com/efritz/api-test/config"
	. "github.com/onsi/gomega"
)

type RequestSuite struct{}

func (s *RequestSuite) TestBuildRequest(t sweet.T) {
	request, body, err := buildRequest(
		&config.Request{
			URL:    testParse("http://test.io/users/{{.UserId}}"),
			Method: "put",
			Headers: map[string][]*tmpl.Template{
				"X-Test": []*tmpl.Template{
					testParse("abcd"),
					testParse("cdef"),
				},
			},
			Body: testParse("Username: {{.Username}}"),
		},
		map[string]interface{}{
			"UserId":   1234,
			"Username": "test",
		},
	)

	Expect(err).To(BeNil())
	Expect(body).To(Equal("Username: test"))
	Expect(request.Method).To(Equal("PUT"))
	Expect(request.URL.String()).To(Equal("http://test.io/users/1234"))
	Expect(request.Header).To(Equal(http.Header{
		"X-Test": []string{"abcd", "cdef"},
	}))
}

func (s *RequestSuite) TestBuildRequestAuth(t sweet.T) {
	request, _, err := buildRequest(
		&config.Request{
			Method: "put",
			Auth: &config.BasicAuth{
				Username: testParse("admin-{{.Index}}"),
				Password: testParse("secret-{{.Index}}"),
			},
		},
		map[string]interface{}{
			"Index": 24,
		},
	)

	Expect(err).To(BeNil())
	username, password, ok := request.BasicAuth()
	Expect(ok).To(BeTrue())
	Expect(username).To(Equal("admin-24"))
	Expect(password).To(Equal("secret-24"))
}

//
// Helpers

func testParse(template string) *tmpl.Template {
	return tmpl.Must(tmpl.New("").Parse(template))
}
