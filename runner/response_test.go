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
		StatusCode: http.StatusAccepted,
		Header: map[string][]string{
			"X-Custom1": []string{"foo"},
			"X-Custom2": []string{"bar", "baz"},
		},
		Body: ioutil.NopCloser(bytes.NewReader([]byte("data: payload"))),
	}

	body, context, errors, err := matchResponse(resp, &config.Response{
		Status: testPattern(`20(\d)`),
		Headers: map[string][]*regexp.Regexp{
			"X-Custom1": []*regexp.Regexp{testPattern(`.+`)},
			"X-Custom2": []*regexp.Regexp{testPattern(`.+`)},
		},
		Body: testPattern(`^data: (.+)`),
	})

	Expect(err).To(BeNil())
	Expect(body).To(Equal("data: payload"))
	Expect(errors).To(BeEmpty())

	Expect(context).To(Equal(map[string]interface{}{
		"statusGroups": []string{"202", "2"},
		"headerGroups": map[string][]string{
			"X-Custom1": []string{"foo"},
			"X-Custom2": []string{"bar"},
		},
		"bodyGroups":       []string{"data: payload", "payload"},
		"extractionGroups": map[string][]string{},
	}))
}

func (s *ResponseSuite) TestMatchResponseMismatchedStatusCode(t sweet.T) {
	resp := &http.Response{
		StatusCode: http.StatusNoContent,
		Body:       ioutil.NopCloser(bytes.NewReader(nil)),
	}

	_, _, errors, err := matchResponse(resp, &config.Response{
		Status: testPattern(`200`),
	})

	Expect(err).To(BeNil())
	Expect(errors).To(ConsistOf(RequestMatchError{
		Type:     "Status Code",
		Expected: "200",
		Actual:   "204",
	}))
}

func (s *ResponseSuite) TestMatchResponseMismatchedHeader(t sweet.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: map[string][]string{
			"X-Custom1": []string{"1235abcd"},
		},
		Body: ioutil.NopCloser(bytes.NewReader(nil)),
	}

	_, _, errors, err := matchResponse(resp, &config.Response{
		Headers: map[string][]*regexp.Regexp{
			"X-Custom1": []*regexp.Regexp{
				testPattern(`\d{8}`),
			},
		},
	})

	Expect(err).To(BeNil())
	Expect(errors).To(ConsistOf(RequestMatchError{
		Type:     "Header 'X-Custom1'",
		Expected: `\d{8}`,
		Actual:   "1235abcd",
	}))
}

func (s *ResponseSuite) TestMatchResponseMismatchedBody(t sweet.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("data: payload"))),
	}

	_, _, errors, err := matchResponse(resp, &config.Response{
		Body: testPattern(`data: (\d+)`),
	})

	Expect(err).To(BeNil())
	Expect(errors).To(ConsistOf(RequestMatchError{
		Type:     "Body",
		Expected: `data: (\d+)`,
		Actual:   "<placeholder>",
	}))
}

// TODO - redo extraction
