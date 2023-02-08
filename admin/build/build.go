package build

import "github.com/urfave/cli/v2"

func GetBuildCmd() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "manage build",
		Subcommands: []*cli.Command{
			getPruneCmd(),
		},
	}
}
