package app

import "github.com/urfave/cli/v3"

const (
	flagLogfile = "log-file"
)

var AppFlags []cli.Flag = []cli.Flag{
	&cli.StringFlag{
		Name:        "color",
		Aliases:     []string{"c"},
		Usage:       "Initial color",
		Value:       "",
		DefaultText: "#b7416e",
	},
	&cli.StringFlag{
		Name:        flagLogfile,
		Aliases:     []string{"l"},
		Usage:       "Log file",
		Sources:     cli.EnvVars("TERMPICKER_LOG_FILE"),
		DefaultText: "/path/to/termpicker-logs.txt",
	},
	cli.VersionFlag,
}
