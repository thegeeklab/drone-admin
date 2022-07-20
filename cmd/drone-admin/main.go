// Copyright (c) 2022, Robert Kaussow <mail@thegeeklab.de>

package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/thegeeklab/drone-admin/admin/build"
	"github.com/urfave/cli/v2"
)

var (
	BuildVersion = "devel"
	BuildDate    = "00000000"
)

func main() {
	if _, err := os.Stat("/run/drone/env"); err == nil {
		_ = godotenv.Overload("/run/drone/env")
	}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version=%s date=%s\n", c.App.Name, c.App.Version, BuildDate)
	}

	app := &cli.App{
		Name:    "drone-admin",
		Usage:   "drone admin tools",
		Version: BuildVersion,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "log level",
				EnvVars: []string{"DRONE_ADMIN_LOG_LEVEL"},
				Value:   "info",
			},
			&cli.StringFlag{
				Name:     "token",
				Aliases:  []string{"t"},
				Usage:    "server auth token",
				EnvVars:  []string{"DRONE_ADMIN_TOKEN"},
				Required: true,
			},

			&cli.StringFlag{
				Name:     "server",
				Aliases:  []string{"s"},
				Usage:    "server address",
				EnvVars:  []string{"DRONE_ADMIN_SERVER"},
				Required: true,
			},
		},
		Before: func(ctx *cli.Context) error {
			lvl, err := logrus.ParseLevel(ctx.String("log-level"))
			if err != nil {
				lvl = logrus.InfoLevel
			}
			logrus.SetLevel(lvl)

			return nil
		},
		Commands: []*cli.Command{
			build.Command,
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
