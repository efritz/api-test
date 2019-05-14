package runner

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/api-test/config"
	. "github.com/onsi/gomega"
)

type ReportSuite struct{}

func (s *ReportSuite) TestFormatJUnitReport(t sweet.T) {
	mockRunner1 := NewMockScenarioRunner()
	mockRunner1.ResultsFunc.SetDefaultReturn([]*TestResult{
		&TestResult{Duration: time.Second * 1},
		&TestResult{Duration: time.Second * 2},
	})

	mockRunner2 := NewMockScenarioRunner()
	mockRunner2.ResultsFunc.SetDefaultReturn([]*TestResult{
		&TestResult{Duration: time.Second * 3},
		&TestResult{Duration: time.Second * 4},
		&TestResult{
			Duration: time.Second * 5,
			RequestMatchErrors: []RequestMatchError{
				RequestMatchError{
					Type:     "Status Code",
					Expected: "2..",
					Actual:   "404",
				},
				RequestMatchError{
					Type:     "Body",
					Expected: "payload",
					Actual:   "not found",
				},
			},
			Request: &http.Request{
				Method:     "POST",
				RequestURI: "/failing-test",
				Header: map[string][]string{
					"X-Req1": []string{"foo"},
					"X-Req2": []string{"bar", "baz"},
				},
			},
			RequestBody: "request payload",
			Response: &http.Response{
				StatusCode: http.StatusNotFound,
				Header: map[string][]string{
					"X-Resp1": []string{"foo"},
					"X-Resp2": []string{"bar", "baz"},
				},
			},
			ResponseBody: "response payload",
		},
	})

	mockRunner3 := NewMockScenarioRunner()
	mockRunner3.ResultsFunc.SetDefaultReturn([]*TestResult{
		&TestResult{Duration: time.Second * 6},
		&TestResult{Duration: time.Second * 7},
		&TestResult{Duration: time.Second * 8, Err: fmt.Errorf("oops")},
	})

	content, err := formatJUnitReport(map[string]*ScenarioContext{
		"s1": &ScenarioContext{
			Scenario: &config.Scenario{
				Name: "s1",
				Tests: []*config.Test{
					&config.Test{Name: "foo1"},
					&config.Test{Name: "bar1"},
					&config.Test{Name: "baz1"},
				},
			},
			Runner: mockRunner1,
		},
		"s2": &ScenarioContext{
			Scenario: &config.Scenario{
				Name: "s2",
				Tests: []*config.Test{
					&config.Test{Name: "foo2"},
					&config.Test{Name: "bar2"},
					&config.Test{Name: "baz2"},
				},
			},
			Runner: mockRunner2,
		},
		"s3": &ScenarioContext{
			Scenario: &config.Scenario{
				Name: "s3",
				Tests: []*config.Test{
					&config.Test{Name: "foo3"},
					&config.Test{Name: "bar3"},
					&config.Test{Name: "baz3"},
				},
			},
			Runner: mockRunner3,
		},
	})

	xmlTemplate := `
	<testsuites>
		<testsuite tests="3" failures="0" time="0.000" name="s1">
			<testcase name="foo1" time="1.000"></testcase>
			<testcase name="bar1" time="2.000"></testcase>
			<testcase name="baz1" time="">
				<skipped message="skipped"></skipped>
			</testcase>
		</testsuite>
		<testsuite tests="3" failures="1" time="0.000" name="s2">
			<testcase name="foo2" time="3.000"></testcase>
			<testcase name="bar2" time="4.000"></testcase>
			<testcase name="baz2" time="5.000">
				<failure type="Assertion Failure" message="Unexpected Status Code, Body">%s</failure>
			</testcase>
		</testsuite>
		<testsuite tests="3" failures="1" time="0.000" name="s3">
			<testcase name="foo3" time="6.000"></testcase>
			<testcase name="bar3" time="7.000"></testcase>
			<testcase name="baz3" time="8.000">
				<failure type="Error" message="oops"></failure>
			</testcase>
		</testsuite>
	</testsuites>
	`

	buffer := bytes.NewBuffer(nil)
	xml.Escape(buffer, []byte(deindent(`
	> Status Code:
		Actual: '404'
		Expected: '2..'

	> Body:
		Actual: 'not found'
		Expected: 'payload'

	> POST <nil>
	> X-Req1: foo
	> X-Req2: bar, baz
	> request payload

	< 404 Not Found
	< X-Resp1: foo
	< X-Resp2: bar, baz
	< response payload
	`)))

	Expect(err).To(BeNil())
	Expect(string(content)).To(MatchXML(fmt.Sprintf(xmlTemplate, buffer.String())))
}

//
// Helpers

func deindent(text string) string {
	parts := strings.Split(text, "\n")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = part[1:]
		}
	}

	return strings.TrimSpace(strings.Join(parts, "\n"))
}
