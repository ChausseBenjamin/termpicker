package main

import "github.com/urfave/cli/v2"

const (
	flagLogfile = "logfile"
)

var AppFlags []cli.Flag = []cli.Flag{
	&cli.StringFlag{
		Name:    flagLogfile,
		Aliases: []string{"l"},
		Usage:   "Log file",
		Value:   "/dev/null", // Don't log by default
	},
	&cli.StringFlag{
		Name:    "color",
		Aliases: []string{"c"},
		Usage:   "Initial color",
		Value:   "",
	},
}
