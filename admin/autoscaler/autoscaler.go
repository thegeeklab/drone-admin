package autoscaler

import "github.com/urfave/cli/v2"

var Command = &cli.Command{
	Name:  "autoscaler",
	Usage: "manage autoscaler",
	Subcommands: []*cli.Command{
		&autoscalerReaperCmd,
	},
}
