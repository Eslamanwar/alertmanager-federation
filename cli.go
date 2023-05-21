package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "none"
)

func run(args []string) error {
	app := cli.NewApp()
	app.Name = "alertmanager-federation-ctl"
	app.Usage = "AlertManager federation cli tool"
	app.Version = fmt.Sprintf("%s (%s)", version, commit)
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "json",
			Usage: "show output as JSON instead of as a table",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "sync-silences",
			Aliases: []string{"sil"},
			Action:  syncSilences,
			Usage:   "Sync silences between alertManagers instances",
			Before:  initializeAPI,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "env",
					Usage:    "env to work on",
					Required: true,
				},
				&cli.IntFlag{
					Name:     "sync-period",
					Value:    60,
					Usage:    "Sync period",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "dry-run",
					Value:    "false",
					Usage:    "Run the command in dry run mode for debugging",
					Required: false,
				},
			},
		},
	}

	err := app.Run(args)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := run(os.Args)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
