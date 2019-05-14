package runner

import (
	"fmt"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ResultSuite struct{}

func (s *ResultSuite) TestErrored(t sweet.T) {
	result1 := &TestResult{Err: fmt.Errorf("oops")}
	result2 := &TestResult{}
	Expect(result1.Errored()).To(BeTrue())
	Expect(result2.Errored()).To(BeFalse())
}

func (s *ResultSuite) TestFailed(t sweet.T) {
	result1 := &TestResult{RequestMatchErrors: []RequestMatchError{RequestMatchError{}}}
	result2 := &TestResult{}
	Expect(result1.Failed()).To(BeTrue())
	Expect(result2.Failed()).To(BeFalse())
}
