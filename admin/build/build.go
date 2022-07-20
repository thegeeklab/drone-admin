package build

import "github.com/urfave/cli/v2"

var Command = &cli.Command{
	Name:  "build",
	Usage: "manage build",
	Subcommands: []*cli.Command{
		&buidPruneCmd,
	},
}
