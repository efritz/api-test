package runner

import (
	"encoding/xml"
	"fmt"
	"time"
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
	suites := JUnitTestSuites{}

	for _, context := range contexts {
		testCases := []JUnitTestCase{}

		for i, test := range context.Scenario.Tests {
			testCase := JUnitTestCase{
				Name:    test.Name,
				Failure: nil,
			}

			if i < len(context.Results) {
				testCase.Time = fmt.Sprintf("%.3f", float64(context.Results[i].Duration)/float64(time.Second))

				if context.Results[i].Failed() {
					testCase.Failure = &JUnitFailure{
						Message:  "Failed",
						Type:     "",
						Contents: "",
					}
				}
			} else {
				testCase.SkipMessage = &JUnitSkipMessage{"skipped"}
			}

			testCases = append(testCases, testCase)
		}

		failures := 0
		for _, result := range context.Results {
			if result.Failed() {
				failures++
			}
		}

		suites.Suites = append(suites.Suites, JUnitTestSuite{
			Tests:     len(context.Scenario.Tests),
			Failures:  failures,
			Time:      fmt.Sprintf("%.3f", float64(context.Duration())/float64(time.Second)),
			Name:      context.Scenario.Name,
			TestCases: testCases,
		})
	}

	return xml.MarshalIndent(suites, "", "\t")
}

// writer := bufio.NewWriter(w)
// 	writer.WriteString(xml.Header)
// writer.Write(bytes)
// writer.WriteByte('\n')
// writer.Flush()
