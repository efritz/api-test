package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/efritz/api-test/loader"
	"github.com/efritz/api-test/logging"
	"github.com/efritz/api-test/runner"
)

type Options struct {
	Colorize        bool
	ConfigPath      string
	Verbose         bool
	JUnitReportPath string
	ForceSequential bool
	Tests           []string // TODO
}

const Version = "0.1.0"

func main() {
	opts := &Options{
		ConfigPath: "api-test.yaml",
	}

	// TODO - whitelist of tests
	app := kingpin.New("api-test", "api-test is a test runner against a local API.").Version(Version)
	app.Flag("color", "Enable colorized output.").Default("true").BoolVar(&opts.Colorize) // TOOD (fix in ij too)
	app.Flag("config", "The path to the config file.").Short('f').StringVar(&opts.ConfigPath)
	app.Flag("verbose", "Output debug logs.").Short('v').Default("false").BoolVar(&opts.Verbose)
	app.Flag("junit", "The path to write a JUnit XML report.").Short('j').StringVar(&opts.JUnitReportPath)
	app.Flag("force-sequential", "Disable parallel execution.").Default("false").BoolVar(&opts.ForceSequential)

	app.Flag("tests", "TODO").StringsVar(&opts.Tests) // TODO - just supply as positional arguments

	if _, err := app.Parse(os.Args[1:]); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}

	logger := logging.NewLogger(
		opts.Colorize,
		opts.Verbose,
	)

	config, err := loader.NewLoader().Load(opts.ConfigPath)
	if err != nil {
		logger.Error("error: %s", err.Error())
		os.Exit(1)
	}

	if opts.ForceSequential {
		// TODO - temporary hack until overrides
		config.Options.ForceSequential = true
	}

	if err := config.EnableTests(opts.Tests); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	runner := runner.NewRunner(
		config,
		logger,
		opts.JUnitReportPath,
	)

	if err := runner.Run(); err != nil {
		logger.Error("error: %s", err.Error())
		os.Exit(1)
	}
}
