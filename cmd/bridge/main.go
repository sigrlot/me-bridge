package bridge

import (
	"fmt"
	"os"

	"github.com/st-chain/me-bridge/cmd/bridge/action"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "bridge",
		Usage: "A blockchain bridge service",
		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "Start the bridge service",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Aliases:  []string{"c"},
						Usage:    "Path to YAML configuration file",
						Required: true,
					},
				},
				Action: action.StartAction,
			},
			{
				Name:  "config",
				Usage: "Configuration management",
				Subcommands: []*cli.Command{
					{
						Name:  "example",
						Usage: "Generate example configuration file",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "output",
								Aliases: []string{"o"},
								Usage:   "Output path for the example config file",
								Value:   "config.yaml",
							},
							&cli.BoolFlag{
								Name:  "force",
								Usage: "Overwrite existing file",
							},
						},
						Action: action.ConfigExampleAction,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
