package jsonconfig

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type AuthSuite struct{}

func (s *AuthSuite) TestTranslate(t sweet.T) {
	auth := &BasicAuth{
		Username: "admin",
		Password: "secret",
	}

	translated, err := auth.Translate()
	Expect(err).To(BeNil())
	Expect(translated).NotTo(BeNil())
	Expect(testExec(translated.Username)).To(Equal("admin"))
	Expect(testExec(translated.Password)).To(Equal("secret"))
}
