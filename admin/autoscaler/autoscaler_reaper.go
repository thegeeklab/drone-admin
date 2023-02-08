package autoscaler

import (
	"os"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/sirupsen/logrus"
	"github.com/thegeeklab/drone-admin/admin/client"
	"github.com/thegeeklab/drone-admin/admin/util"
	"github.com/urfave/cli/v2"
)

func getReaperCmd() *cli.Command {
	return &cli.Command{
		Name:   "reaper",
		Usage:  "find and kill agents in error state",
		Action: reaper,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "state-file",
				Usage:   "state file",
				EnvVars: []string{"DRONE_ADMIN_AUTOSCALER_REAPER_STATE_FILE"},
				Value:   "/tmp/droneclean.gob",
			},
		},
	}
}

func reaper(ctx *cli.Context) error {
	const maxRetries = 3

	statefile := ctx.String("state-file")
	scaler := ctx.StringSlice("server")
	dry := ctx.Bool("dry-run")
	state := map[string]int{}
	force := false

	if dry {
		logrus.Info("dry-run enabled, no data will be removed")
	}

	if _, err := os.Stat(statefile); err == nil {
		err = util.ReadGob(statefile, &state)
		if err != nil {
			return err
		}
	}

	for _, scaler := range scaler {
		client, err := client.New(scaler, ctx.String("token"))
		if err != nil {
			return err
		}

		servers, err := getServers(client)
		if err != nil {
			return err
		}

		serversAll := len(servers)
		servers = util.Filter(servers, func(s *drone.Server) bool {
			return s.State == "running"
		})

		searchFields := logrus.Fields{
			"server": scaler,
			"ok":     serversAll,
			"error":  len(servers),
		}
		logrus.WithFields(searchFields).Infof("lookup agents in error state")

		for _, server := range servers {
			state[server.Name]++
			triage := state[server.Name]

			if state[server.Name] == maxRetries {
				force = true

				delete(state, server.Name)
			}

			foundFields := logrus.Fields{
				"server": scaler,
				"agent":  server.Name,
				"triage": triage,
				"force":  force,
			}
			logrus.WithFields(foundFields).Infof("destroy agent")

			if !dry {
				err = serverDestroy(client, server.Name, force)
				if err != nil && !strings.Contains(err.Error(), "client error 404") {
					return err
				}
			}
		}
	}

	if !dry {
		err := util.WriteGob(statefile, state)
		if err != nil {
			return err
		}
	}

	return nil
}

func getServers(client drone.Client) ([]*drone.Server, error) {
	servers, err := client.ServerList()
	if err != nil {
		return nil, err
	}

	servers = util.Filter(servers, func(s *drone.Server) bool {
		return s.State != "stopped"
	})

	return servers, nil
}

func serverDestroy(client drone.Client, server string, force bool) error {
	err := client.ServerDelete(server, force)
	if err != nil {
		return err
	}

	return nil
}
