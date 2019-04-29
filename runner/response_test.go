package runner

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/aphistic/sweet"
	"github.com/efritz/api-test/config"
	. "github.com/onsi/gomega"
)

type ResponseSuite struct{}

func (s *ResponseSuite) TestMatchResponse(t sweet.T) {
	resp := &http.Response{
		// TODO
		Body: ioutil.NopCloser(bytes.NewReader([]byte("payload"))),
	}

	body, context, errors, err := matchResponse(resp, &config.Response{
		// TODO
	})

	Expect(err).To(BeNil())
	Expect(body).To(Equal("payload"))
	Expect(errors).To(BeEmpty())
	Expect(context).To(Equal(map[string]interface{}{
		// TODO
	}))
}

//
// Helpers

func testCompile(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}
