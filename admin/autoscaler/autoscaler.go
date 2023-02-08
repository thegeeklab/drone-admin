package autoscaler

import "github.com/urfave/cli/v2"

func GetAutoscalerCmd() *cli.Command {
	return &cli.Command{
		Name:  "autoscaler",
		Usage: "manage autoscaler",
		Subcommands: []*cli.Command{
			getReaperCmd(),
		},
	}
}
