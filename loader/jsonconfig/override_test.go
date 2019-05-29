package jsonconfig

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type OverrideSuite struct{}

func (s *OverrideSuite) TestTranslate(t sweet.T) {
	options := &Override{
		Options: &Options{
			ForceSequential: true,
		},
	}

	translated, err := options.Translate()
	Expect(err).To(BeNil())
	Expect(translated.Options.ForceSequential).To(BeTrue())
}
