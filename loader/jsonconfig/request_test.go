package jsonconfig

import (
	"encoding/json"
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type RequestSuite struct{}

func (s *RequestSuite) TestTranslate(t sweet.T) {
	// TODO
}

func (s *RequestSuite) TestTranslateWithJSONBody(t sweet.T) {
	// TODO
}

func (s *RequestSuite) TestTranslateAbsoluteURI(t sweet.T) {
	// TODO
}

func (s *RequestSuite) TestTranslateNoExplicitMethod(t sweet.T) {
	// TODO
}

func (s *RequestSuite) TestTranslateStringLists(t sweet.T) {
	// TODO
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
	Expect(err).To(MatchError("illegal uri template"))
}

func (s *RequestSuite) TestTranslateInvalidHeaderTemplate(t sweet.T) {
	request := &Request{
		Headers: map[string]json.RawMessage{
			"X-Request-ID": json.RawMessage(`"{{"`),
		},
	}

	_, err := request.Translate(nil)
	Expect(err).To(MatchError("illegal header template"))
}

func (s *RequestSuite) TestTranslateInvalidBodyTemplate(t sweet.T) {
	request := &Request{
		Body: "{{",
	}

	_, err := request.Translate(nil)
	Expect(err).To(MatchError("illegal body template"))
}

func (s *RequestSuite) TestTranslateInvalidJSONBodyTemplate(t sweet.T) {
	request := &Request{
		JSONBody: json.RawMessage(`{"foo": "{{"}`),
	}

	_, err := request.Translate(nil)
	Expect(err).To(MatchError("illegal json body template"))
}
