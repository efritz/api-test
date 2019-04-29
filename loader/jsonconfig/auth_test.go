package jsonconfig

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/efritz/api-test/config"
)

type AuthSuite struct{}

func (s *AuthSuite) TestTranslate(t sweet.T) {
	auth := &BasicAuth{
		Username: "admin",
		Password: "secret",
	}

	translated, err := auth.Translate()
	Expect(err).To(BeNil())
	Expect(translated).To(Equal(&config.BasicAuth{
		Username: "admin",
		Password: "secret",
	}))
}
