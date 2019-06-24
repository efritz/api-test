package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/hokaccha/go-prettyjson"

	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/logging"
	"github.com/efritz/pentimento"
)

func displayProgress(
	logger logging.Logger,
	names []string,
	contexts map[string]*ScenarioContext,
	p *pentimento.Printer,
) {
	content := pentimento.NewContent()

	for _, name := range names {
		context := contexts[name]

		if context.Scenario.Disabled {
			continue
		}

		details := ""
		for _, result := range context.Runner.Results() {
			if result == nil {
				continue
			}

			if result.Errored() {
				details += logger.Colorize(logging.ColorError, "E")
			} else if result.Failed() {
				details += logger.Colorize(logging.ColorError, "F")
			} else if result.Skipped {
				details += logger.Colorize(logging.ColorWarn, "S")
			} else {
				details += logger.Colorize(logging.ColorInfo, ".")
			}
		}

		if context.Runner.Resolved() {
			details += fmt.Sprintf(
				" (in %s)",
				formatMilliseconds(context.Runner.Duration()),
			)
		}

		content.AddLine(
			"[%s] Scenario %s %s",
			getStatus(logger, context),
			name,
			details,
		)
	}

	p.WriteContent(content)
}

func displaySummary(
	logger logging.Logger,
	contexts map[string]*ScenarioContext,
	started time.Time,
) {
	wallDuration := time.Now().Sub(started)

	totalDuration := time.Duration(0)
	for _, context := range contexts {
		totalDuration += context.Runner.Duration()
	}

	numScenarios := 0
	numScenariosSkipped := 0
	numTests := 0
	numTestsSkipped := 0
	numFailures := 0

	for _, context := range contexts {
		if context.Scenario.Disabled {
			continue
		}

		if !context.Skipped {
			numScenarios++
		} else {
			numScenariosSkipped++
		}

		for _, result := range context.Runner.Results() {
			if result == nil {
				continue
			}

			if !result.Skipped {
				numTests++
			} else {
				numTestsSkipped++
			}
		}

		if context.Runner.Errored() || context.Runner.Failed() {
			numFailures++
		}
	}

	logger.Log(nil, "")

	if numScenariosSkipped > 0 || numTestsSkipped > 0 {
		logger.Log(
			nil,
			logger.Colorize(
				logging.ColorWarn,
				"Skipped %d scenarios and %d tests",
				numScenariosSkipped,
				numTestsSkipped,
			),
		)
	}

	if numFailures == 0 {
		logger.Log(
			nil,

			logger.Colorize(
				logging.ColorInfo,
				"Ran %d scenarios and %d tests in %s (%s on the wall)",
				numScenarios,
				numTests,
				formatSeconds(totalDuration),
				formatSeconds(wallDuration),
			),
		)

		return
	}

	logger.Log(
		nil,
		logger.Colorize(
			logging.ColorError,
			"Failed %d out of %d ran\n",
			numFailures,
			numScenarios,
		),
	)

	for _, context := range contexts {
		for i, result := range context.Runner.Results() {
			if result == nil || (!result.Errored() && !result.Failed()) {
				continue
			}

			logger.Log(
				nil,
				logger.Colorize(
					logging.ColorError,
					"%s/%s: ",
					context.Scenario.Name,
					context.Scenario.Tests[i].Name,
				),
			)

			displayFailure(
				logger,
				context.Scenario,
				context.Scenario.Tests[i],
				result,
			)
		}
	}
}

func displayFailure(
	logger logging.Logger,
	scenario *config.Scenario,
	test *config.Test,
	result *TestResult,
) {
	if result.Err != nil {
		logger.Log(
			nil,
			logger.Colorize(
				logging.ColorError,
				"Failed to perform request: %s",
				result.Err.Error(),
			),
		)

		return
	}

	for _, err := range result.RequestMatchErrors {
		logger.Log(
			nil,
			logger.Colorize(
				logging.ColorError,
				"> %s:\n\tActual: '%s'\n\tExpected: '%s'",
				err.Type,
				err.Actual,
				err.Expected,
			),
		)

		logger.Log(nil, "")
	}

	logger.Log(nil, formatRequest(result.Request, result.RequestBody, logger.Colorized()))
	logger.Log(nil, formatResponse(result.Response, result.ResponseBody, logger.Colorized()))
}

func getStatus(
	logger logging.Logger,
	context *ScenarioContext,
) *pentimento.AnimatedString {
	statuses := map[bool]*pentimento.AnimatedString{
		true:                      pentimento.ScrollingDots,
		context.Pending:           pentimento.NewStaticString("   "),
		context.Skipped:           pentimento.NewStaticString(logger.Colorize(logging.ColorWarn, " ✗ ")),
		context.Runner.Errored():  pentimento.NewStaticString(logger.Colorize(logging.ColorError, " ✗ ")),
		context.Runner.Failed():   pentimento.NewStaticString(logger.Colorize(logging.ColorError, " ✗ ")),
		context.Runner.Resolved(): pentimento.NewStaticString(logger.Colorize(logging.ColorInfo, " ✓ ")),
	}

	return statuses[true]
}

func formatMilliseconds(duration time.Duration) string {
	if duration < time.Second {
		return fmt.Sprintf("%dms", int(duration)/int(time.Millisecond))
	}

	return fmt.Sprintf("%.2fs", float64(duration)/float64(time.Second))
}

func formatSeconds(duration time.Duration) string {
	if duration < time.Minute {
		return fmt.Sprintf("%ds", int(duration)/int(time.Second))
	}

	return fmt.Sprintf("%.2fm", float64(duration)/float64(time.Minute))
}

func formatRequest(req *http.Request, body string, colorize bool) string {
	line := fmt.Sprintf(
		"%s %s\n%s\n%s\n",
		strings.ToUpper(req.Method),
		req.URL,
		formatHeaders(req.Header),
		formatBody(body, req.Header, colorize),
	)

	return fmt.Sprintf("%s\n", prefix(">", line))
}

func formatResponse(resp *http.Response, body string, colorize bool) string {
	line := fmt.Sprintf(
		"%d %s\n%s\n%s\n",
		resp.StatusCode,
		http.StatusText(resp.StatusCode),
		formatHeaders(resp.Header),
		formatBody(body, resp.Header, colorize),
	)

	return fmt.Sprintf("%s\n", prefix("<", line))
}

func formatBody(body string, headers http.Header, colorize bool) string {
	if headers.Get("Content-Type") == "application/json" {
		if colorize {
			if s, err := prettyjson.Marshal(json.RawMessage(body)); err == nil {
				return string(s)
			}
		}

		out := bytes.Buffer{}
		json.Indent(&out, []byte(body), "", "  ")
		return out.String()
	}

	return body
}

func formatHeaders(headers http.Header) string {
	lines := []string{}
	for key, values := range headers {
		lines = append(lines, fmt.Sprintf("%s: %s", key, strings.Join(values, ", ")))
	}

	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

func prefix(prefix, text string) string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	for i, line := range lines {
		lines[i] = fmt.Sprintf("%s %s", prefix, line)
	}

	return strings.Join(lines, "\n")
}
