package app

import "github.com/urfave/cli/v3"

const (
	flagLogfile   = "log-file"
	flagColor     = "color"
	flagSampleStr = "sample-text"
	flagSampleBG  = "background-sample"
	flagSampleFG  = "foreground-sample"
)

var AppFlags []cli.Flag = []cli.Flag{
	&cli.StringFlag{
		Name:        flagColor,
		Usage:       "Initial color",
		Aliases:     []string{"c"},
		Value:       "",
		DefaultText: "#b7416e",
	},
	&cli.StringFlag{
		Name:        flagLogfile,
		Usage:       "Log file",
		Aliases:     []string{"l"},
		Sources:     cli.EnvVars("TERMPICKER_LOG_FILE"),
		DefaultText: "/path/to/termpicker-logs.txt",
	},
	&cli.StringFlag{
		Name:    flagSampleStr,
		Usage:   "Text to preview colors as a foreground/background",
		Sources: cli.EnvVars("TERMPICKER_PREVIEW_STRING"),
		Aliases: []string{"t"},
		Value:   "The quick brown fox jump over the lazy dog",
	},
	&cli.StringFlag{
		Name:        flagSampleBG,
		Usage:       "Color used in the background when previewing target color as text/foreground (requires `sample-text` to be set)",
		Aliases:     []string{"bg"},
		DefaultText: "#111a1f",
		Value:       "",
	},
	&cli.StringFlag{
		Name:        flagSampleFG,
		Aliases:     []string{"fg"},
		Usage:       "Color used for foreground/text when previewing target color as a background (requires `sample-text` to be set)",
		DefaultText: "#ebcb88",
		Value:       "",
	},
	cli.VersionFlag,
}
