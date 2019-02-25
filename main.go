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
	Colorize   bool
	ConfigPath string
	Quiet      bool
	Verbose    bool
}

const Version = "0.1.0"

func main() {
	opts := &Options{
		ConfigPath: "test.yaml",
	}

	app := kingpin.New("api-test", "api-test is a test runner against a local API.").Version(Version)
	app.Flag("color", "Enable colorized output.").Default("true").BoolVar(&opts.Colorize)
	app.Flag("config", "The path to the config file.").Short('f').StringVar(&opts.ConfigPath)
	app.Flag("quiet", "Do not output to stdout or stderr.").Short('q').Default("false").BoolVar(&opts.Quiet)
	app.Flag("verbose", "Output debug logs.").Short('v').Default("false").BoolVar(&opts.Verbose)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}

	logger := logging.NewLogger(
		opts.Colorize,
		opts.Quiet,
		opts.Verbose,
	)

	config, err := loader.Load(opts.ConfigPath)
	if err != nil {
		logger.Error("error: %s", err.Error())
		os.Exit(1)
	}

	runner := runner.NewRunner(config, logger)

	if err := runner.Run(); err != nil {
		logger.Error("error: %s", err.Error())
		os.Exit(1)
	}
}
