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
	Colorize        bool
	ConfigPath      string
	JUnitReportPath string
	ForceSequential bool
}

const Version = "0.1.0"

func main() {
	opts := &Options{}
	app := kingpin.New("api-test", "api-test is a test runner against a local API.").Version(Version)
	app.Flag("color", "Enable colorized output.").Default("true").BoolVar(&opts.Colorize) // TOOD (fix in ij too)
	app.Flag("config", "The path to the config file.").Short('f').StringVar(&opts.ConfigPath)
	app.Flag("junit", "The path to write a JUnit XML report.").Short('j').StringVar(&opts.JUnitReportPath)
	app.Flag("force-sequential", "Disable parallel execution.").Default("false").BoolVar(&opts.ForceSequential)
	tests := app.Arg("tests", "A list of specific scenarios or tests to run.").Strings()

	if _, err := app.Parse(os.Args[1:]); err != nil {
		logging.EmergencyLog("error: %s", err.Error())
		os.Exit(1)
	}

	logger := logging.NewLogger(
		opts.Colorize,
	)

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
		runner.WithJUnitReportPath(opts.JUnitReportPath),
	)

	if err := runner.Run(); err != nil {
		logger.Error("error: %s", err.Error())
		os.Exit(1)
	}
}
