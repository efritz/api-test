package main

import (
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/efritz/api-test/config"
	"github.com/efritz/api-test/loader"
	"github.com/efritz/api-test/logging"
	"github.com/efritz/api-test/runner"
)

type Options struct {
	ConfigPath      string
	Env             []string
	DisableColor    bool
	JUnitReportPath string
	ForceSequential bool
	Verbose         bool
}

const Version = "0.1.0"

func main() {
	opts := &Options{}
	app := kingpin.New("api-test", "api-test is a test runner against a local API.").Version(Version)
	app.Flag("config", "The path to the config file.").Short('f').StringVar(&opts.ConfigPath)
	app.Flag("env", "Environment variables.").Short('e').StringsVar(&opts.Env)
	app.Flag("no-color", "Disable colorized output.").Default("false").BoolVar(&opts.DisableColor)
	app.Flag("junit", "The path to write a JUnit XML report.").Short('j').StringVar(&opts.JUnitReportPath)
	app.Flag("force-sequential", "Disable parallel execution.").Default("false").BoolVar(&opts.ForceSequential)
	app.Flag("verbose", "Enable verbose logging.").Short('v').BoolVar(&opts.Verbose)
	tests := app.Arg("tests", "A list of specific scenarios or tests to run.").Strings()

	if _, err := app.Parse(os.Args[1:]); err != nil {
		logging.EmergencyLog("error: %s", err.Error())
		os.Exit(1)
	}

	logger := logging.NewLogger(!opts.DisableColor)

	path, err := loader.GetConfigPath(opts.ConfigPath)
	if err != nil {
		logger.Error("error: %s", err.Error())
		os.Exit(1)
	}

	override := &config.Override{
		Options: &config.Options{
			ForceSequential: opts.ForceSequential,
		},
	}

	config, err := loader.Load(path, override)
	if err != nil {
		logger.Error("error: %s", err.Error())
		os.Exit(1)
	}

	if err := config.EnableTests(*tests); err != nil {
		logger.Error("error: %s", err.Error())
		os.Exit(1)
	}

	runner := runner.NewRunner(
		config,
		runner.WithLogger(logger),
		runner.WithEnvironment(opts.Env),
		runner.WithJUnitReportPath(opts.JUnitReportPath),
		runner.WithVerboseLogging(opts.Verbose),
	)

	if err := runner.Run(); err != nil {
		logger.Error("error: %s", err.Error())
		os.Exit(1)
	}
}
