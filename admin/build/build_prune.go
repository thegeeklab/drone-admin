package build

import (
	"sort"
	"time"

	"github.com/drone/drone-go/drone"
	"github.com/sirupsen/logrus"
	"github.com/thegeeklab/drone-admin/admin/client"
	"github.com/thegeeklab/drone-admin/admin/util"
	"github.com/urfave/cli/v2"
)

var buidPruneCmd = cli.Command{
	Name:      "prune",
	Usage:     "prune builds",
	ArgsUsage: "<repo/name>",
	Action:    buidPrune,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "older-than",
			Usage:    "remove builds older than the specified time limit",
			EnvVars:  []string{"DRONE_ADMIN_BUILD_PRUNE_OLDER_THAN"},
			Required: true,
		},
		&cli.IntFlag{
			Name:    "keep-min",
			Usage:   "minimum number of builds to keep",
			EnvVars: []string{"DRONE_ADMIN_BUILD_PRUNE_KEEP_MIN"},
			Value:   10,
		},
		&cli.BoolFlag{
			Name:    "dry-run",
			Usage:   "disable api calls",
			EnvVars: []string{"DRONE_ADMIN_BUILD_PRUNE_DRY_RUN"},
			Value:   false,
		},
	},
}

func buidPrune(c *cli.Context) error {
	client, err := client.New(c.String("server"), c.String("token"))
	if err != nil {
		return err
	}

	repos, err := getRepos(client)
	if err != nil {
		return err
	}

	dur := c.String("older-than")
	keep := c.Int("keep-min")
	dry := c.Bool("dry-run")

	duration, err := time.ParseDuration(dur)
	if err != nil {
		return err
	}

	if dry {
		logrus.Info("dry-run enabled, no data will be removed")
	}

	logrus.Infof("prune builds older than %v, keep min %v", dur, keep)

	for _, r := range repos {
		builds, err := getBuilds(client, r)
		if err != nil {
			return err
		}

		sort.Slice(builds, func(i, j int) bool {
			return builds[i].Number > builds[j].Number
		})

		if bl := len(builds); bl > 0 && bl > keep {
			builds = builds[keep:]
			builds = util.Filter(builds, func(b *drone.Build) bool {
				return time.Since(time.Unix(b.Created, 0)) > duration
			})

			logrus.Infof("prune %v/%v builds from '%v'", len(builds), bl, r.Slug)

			if !dry && len(builds) > 0 {
				err := client.BuildPurge(r.Namespace, r.Name, int(builds[0].Number+1))
				if err != nil {
					return err
				}
			}
		} else {
			logrus.Infof("skip '%v', number of %v builds lower than min value", r.Slug, len(builds))
		}

	}

	return nil
}

func getRepos(client drone.Client) ([]*drone.Repo, error) {
	page := 1
	repos := make([]*drone.Repo, 0)

	for {
		r, err := client.RepoListAll(drone.ListOptions{Page: page, Size: 50})
		if err != nil {
			return nil, err
		}

		if len(r) == 0 {
			break
		}

		repos = append(repos, r...)
		page++
	}

	repos = util.Filter(repos, func(r *drone.Repo) bool {
		return r.Active
	})

	return repos, nil
}

func getBuilds(client drone.Client, repo *drone.Repo) ([]*drone.Build, error) {
	page := 1
	builds := make([]*drone.Build, 0)

	for {
		b, err := client.BuildList(repo.Namespace, repo.Name, drone.ListOptions{Page: page, Size: 50})
		if err != nil {
			return nil, err
		}

		if len(b) == 0 {
			break
		}

		builds = append(builds, b...)
		page++
	}
	return builds, nil
}
