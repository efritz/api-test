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

		details := ""
		if !context.Pending && !context.Skipped && !context.Resolved() && !context.Failed() {
			details = fmt.Sprintf(
				"(%d/%d)",
				len(context.Results)+1,
				len(context.Scenario.Tests),
			)
		}

		if context.Resolved() {
			details = fmt.Sprintf(
				"finished in %s",
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
	numTests := 0
	numFailures := 0

	for _, context := range contexts {
		if !context.Skipped {
			numScenarios++
		}

		numTests += len(context.Results)

		if context.Failed() {
			numFailures++
		}
	}

	if numFailures == 0 {
		logger.Info(
			"\nRan %d scenarios and %d tests in %s (%s on the wall)",
			numScenarios,
			numTests,
			formatSeconds(totalDuration),
			formatSeconds(wallDuration),
		)

		return
	}

	logger.Error(
		"\nFailed %d out of %d ran\n",
		numFailures,
		numScenarios,
	)

	for _, context := range contexts {
		if !context.Failed() {
			continue
		}

		if lastResult := context.LastResult(); lastResult != nil {
			displayFailure(context.Scenario, context.LastTest(), lastResult, logger)
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

	if lastResult := context.LastResult(); lastResult != nil {
		if context.Resolved() {
			return passStatus
		}

		if len(lastResult.RequestMatchErrors) > 0 {
			return failedStatus
		}

		if lastResult.Err != nil {
			return errorStatus
		}
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
