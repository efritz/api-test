package runner

import (
	"bytes"
	"io/ioutil"
	"net/http"

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
		Body: ioutil.NopCloser(bytes.NewReader([]byte(`{
			"data": "payload"
		}`))),
	}

	body, context, errors, err := matchResponse(resp, &config.Response{
		Status: testPattern(`20(\d)`),
		Extract: map[string]*config.ValueExtractor{
			"foo": &config.ValueExtractor{
				JQ: ".data",
				Assert: &config.ValueAssertion{
					Pattern: testPattern("pay(.+)"),
				},
			},
			"bar": &config.ValueExtractor{
				Pattern: testPattern(`.+`),
				Header:  "X-Custom1",
			},
			"baz": &config.ValueExtractor{
				Pattern: testPattern(`.+`),
				Header:  "X-Custom2",
			},
		},
	})

	Expect(err).To(BeNil())
	Expect(body).To(MatchJSON(`{"data": "payload"}`))
	Expect(errors).To(BeEmpty())

	Expect(context).To(Equal(map[string]interface{}{
		"foo":    "payload",
		"bar":    []string{"foo"},
		"baz":    []string{"bar"},
		"status": 202,
		"headers": map[string]string{
			"X-Custom1": "foo",
			"X-Custom2": "bar",
		},
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

func (s *ResponseSuite) TestMatchResponseJQExtractFailure(t sweet.T) {
	resp := &http.Response{
		StatusCode: http.StatusNoContent,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{invalid`))),
	}

	_, _, errors, err := matchResponse(resp, &config.Response{
		Status: testPattern(`2..`),
		Extract: map[string]*config.ValueExtractor{
			"foo": &config.ValueExtractor{
				JQ: ".data[].foo",
				Assert: &config.ValueAssertion{
					Pattern: testPattern(`\d{4}`),
				},
			},
		},
	})

	Expect(err).To(BeNil())
	Expect(errors).To(ConsistOf(RequestMatchError{
		Type:     "Body",
		Expected: ".data[].foo",
		Actual:   "{invalid",
	}))
}

func (s *ResponseSuite) TestMatchResponseRegexExtractFailure(t sweet.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: map[string][]string{
			"X-Custom1": []string{"1235abcd"},
		},
		Body: ioutil.NopCloser(bytes.NewReader(nil)),
	}

	_, _, errors, err := matchResponse(resp, &config.Response{
		Extract: map[string]*config.ValueExtractor{
			"foo": &config.ValueExtractor{
				Pattern: testPattern(`\d{8}`),
				Header:  "X-Custom1",
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

func (s *ResponseSuite) TestMatchResponseRegexAssertionFailure(t sweet.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: ioutil.NopCloser(bytes.NewReader([]byte(`{
			"data": [
				{"foo": 123},
				{"foo": 234},
				{"foo": 345}
			]
		}`))),
	}

	_, _, errors, err := matchResponse(resp, &config.Response{
		Extract: map[string]*config.ValueExtractor{
			"foo": &config.ValueExtractor{
				JQ: ".data[].foo",
				Assert: &config.ValueAssertion{
					Pattern: testPattern(`\d{4}`),
				},
			},
		},
	})

	Expect(err).To(BeNil())
	Expect(errors).To(ConsistOf(RequestMatchError{
		Type:     "Body",
		Expected: `\d{4}`,
		Actual:   "123",
	}))
}

func (s *ResponseSuite) TestMatchResponseJSONSchemaAssertionFailure(t sweet.T) {
	data := `{
		"data": [
			{"foo": 123},
			{"foo": 234},
			{"foo": 345}
		]
	}`

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(data))),
	}

	_, _, errors, err := matchResponse(resp, &config.Response{
		Extract: map[string]*config.ValueExtractor{
			"foo": &config.ValueExtractor{
				JQ: ".",
				Assert: &config.ValueAssertion{
					Schema: testSchema(`{
						"type": "array",
						"item": {
							"type": "string"
						}
					}`),
				},
			},
		},
	})

	Expect(err).To(BeNil())
	Expect(errors[0].Actual).To(MatchJSON(data))
}
