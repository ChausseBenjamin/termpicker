package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ChausseBenjamin/termpicker/internal/app"
	"github.com/ChausseBenjamin/termpicker/internal/util"
)

var version = "compiled"

func main() {
	app := app.Command(version)
	if err := app.Run(context.Background(), os.Args); err != nil {
		slog.Error("Program crashed", util.ErrKey, err.Error())
		os.Exit(1)
	}
}
