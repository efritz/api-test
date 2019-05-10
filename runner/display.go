package runner

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/logging"
	"github.com/efritz/pentimento"
)

var (
	runningStatus = pentimento.ScrollingDots
	pendingStatus = pentimento.NewStaticString("   ")
	skippedStatus = pentimento.NewStaticString(logging.Colorize(" ✗ ", logging.LevelWarn))
	passStatus    = pentimento.NewStaticString(logging.Colorize(" ✓ ", logging.LevelInfo))
	failedStatus  = pentimento.NewStaticString(logging.Colorize(" ✗ ", logging.LevelError))
	errorStatus   = pentimento.NewStaticString(logging.Colorize(" ✗ ", logging.LevelError))
)

func displayProgress(
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
		for _, result := range context.Results {
			if result == nil {
				continue
			}

			if result.Errored() {
				details += logging.Colorize("E", logging.LevelError)
			} else if result.Failed() {
				details += logging.Colorize("F", logging.LevelError)
			} else if result.Skipped {
				details += logging.Colorize("S", logging.LevelWarn)
			} else {
				details += logging.Colorize(".", logging.LevelInfo)
			}
		}

		if context.Resolved() {
			details += fmt.Sprintf(
				" (in %s)",
				formatMilliseconds(context.Duration()),
			)
		}

		content.AddLine(
			"[%s] Scenario %s %s",
			getStatus(context),
			name,
			details,
		)
	}

	p.WriteContent(content)
}

func displaySummary(
	contexts map[string]*ScenarioContext,
	started time.Time,
	logger logging.Logger,
) {
	wallDuration := time.Now().Sub(started)

	totalDuration := time.Duration(0)
	for _, context := range contexts {
		totalDuration += context.Duration()
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

		if context.Skipped {
			numScenariosSkipped++
		}

		if !context.Skipped {
			numScenarios++
		}

		for _, result := range context.Results {
			if result == nil {
				continue
			}

			if !result.Skipped {
				numTests++
			} else {
				numTestsSkipped++
			}
		}

		if context.Errored() || context.Failed() {
			numFailures++
		}
	}

	logger.Info("")

	if numScenariosSkipped > 0 || numTestsSkipped > 0 {
		logger.Warn(
			"Skipped %d scenarios and %d tests",
			numScenariosSkipped,
			numTestsSkipped,
		)
	}

	if numFailures == 0 {
		logger.Info(
			"Ran %d scenarios and %d tests in %s (%s on the wall)",
			numScenarios,
			numTests,
			formatSeconds(totalDuration),
			formatSeconds(wallDuration),
		)

		return
	}

	logger.Error(
		"Failed %d out of %d ran\n",
		numFailures,
		numScenarios,
	)

	for _, context := range contexts {
		for i, result := range context.Results {
			if result == nil || (!result.Errored() && !result.Failed()) {
				continue
			}

			displayFailure(
				context.Scenario,
				context.Scenario.Tests[i],
				result,
				logger,
			)
		}
	}
}

func displayFailure(
	scenario *config.Scenario,
	test *config.Test,
	result *TestResult,
	logger logging.Logger,
) {
	logger.Error("%s/%s: ", scenario.Name, test.Name)

	if result.Err != nil {
		logger.Error("Failed to perform request: %s", result.Err.Error())
		return
	}

	for _, err := range result.RequestMatchErrors {
		logger.Error(
			"> %s:\n\t  Actual: '%s'\n\tExpected: '%s'",
			err.Type,
			err.Actual,
			err.Expected,
		)

		logger.Error("")
	}

	logger.Info(formatRequest(result.Request, result.RequestBody))
	logger.Info(formatResponse(result.Response, result.ResponseBody))
}

func getStatus(context *ScenarioContext) *pentimento.AnimatedString {
	if context.Running {
		return runningStatus
	}

	if context.Skipped {
		return skippedStatus
	}

	if context.Resolved() {
		return passStatus
	}

	if context.Errored() {
		return errorStatus
	}

	if context.Failed() {
		return failedStatus
	}

	return pendingStatus
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

func formatRequest(req *http.Request, body string) string {
	line := fmt.Sprintf(
		"%s %s\n%s\n%s\n",
		strings.ToUpper(req.Method),
		req.URL,
		formatHeaders(req.Header),
		body,
	)

	return fmt.Sprintf("%s\n", prefix(">", line))
}

func formatResponse(resp *http.Response, body string) string {
	line := fmt.Sprintf(
		"%d %s\n%s\n%s\n",
		resp.StatusCode,
		http.StatusText(resp.StatusCode),
		formatHeaders(resp.Header),
		body,
	)

	return fmt.Sprintf("%s\n", prefix("<", line))
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
