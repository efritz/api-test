package runner

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/efritz/api-test/logging"
)

type (
	JUnitTestSuites struct {
		XMLName xml.Name `xml:"testsuites"`
		Suites  []JUnitTestSuite
	}

	JUnitTestSuite struct {
		XMLName   xml.Name `xml:"testsuite"`
		Tests     int      `xml:"tests,attr"`
		Failures  int      `xml:"failures,attr"`
		Time      string   `xml:"time,attr"`
		Name      string   `xml:"name,attr"`
		TestCases []JUnitTestCase
	}

	JUnitTestCase struct {
		XMLName     xml.Name          `xml:"testcase"`
		Name        string            `xml:"name,attr"`
		Time        string            `xml:"time,attr"`
		Failure     *JUnitFailure     `xml:"failure,omitempty"`
		SkipMessage *JUnitSkipMessage `xml:"skipped,omitempty"`
	}

	JUnitFailure struct {
		Message  string `xml:"message,attr"`
		Type     string `xml:"type,attr"`
		Contents string `xml:",chardata"`
	}

	JUnitSkipMessage struct {
		Message string `xml:"message,attr"`
	}
)

func formatJUnitReport(contexts map[string]*ScenarioContext) ([]byte, error) {
	names := []string{}
	for name := range contexts {
		names = append(names, name)
	}

	sort.Strings(names)

	suites := JUnitTestSuites{}
	for _, name := range names {
		context := contexts[name]
		testCases := []JUnitTestCase{}
		results := context.Runner.Results()

		for i, test := range context.Scenario.Tests {
			testCase := JUnitTestCase{
				Name:    test.Name,
				Failure: nil,
			}

			if i < len(results) {
				testCase.Time = fmt.Sprintf("%.3f", float64(results[i].Duration)/float64(time.Second))

				if results[i].Errored() {
					testCase.Failure = &JUnitFailure{
						Type:    "Error",
						Message: results[i].Err.Error(),
					}
				}

				if results[i].Failed() {
					types := []string{}
					for _, err := range results[i].RequestMatchErrors {
						types = append(types, err.Type)
					}

					logger := logging.NewStringLogger()

					displayFailure(
						logger,
						context.Scenario,
						test,
						results[i],
					)

					testCase.Failure = &JUnitFailure{
						Type:     "Assertion Failure",
						Message:  fmt.Sprintf("Unexpected %s", strings.Join(types, ", ")),
						Contents: strings.TrimSpace(logger.String()),
					}
				}
			} else {
				testCase.SkipMessage = &JUnitSkipMessage{"skipped"}
			}

			testCases = append(testCases, testCase)
		}

		failures := 0
		for _, result := range results {
			if result.Errored() || result.Failed() {
				failures++
			}
		}

		suites.Suites = append(suites.Suites, JUnitTestSuite{
			Tests:     len(context.Scenario.Tests),
			Failures:  failures,
			Time:      fmt.Sprintf("%.3f", float64(context.Runner.Duration())/float64(time.Second)),
			Name:      context.Scenario.Name,
			TestCases: testCases,
		})
	}

	return xml.MarshalIndent(suites, "", "\t")
}
