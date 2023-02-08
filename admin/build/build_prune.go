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

const DroneClientPageSize = 50

func getPruneCmd() *cli.Command {
	return &cli.Command{
		Name:      "prune",
		Usage:     "prune builds",
		ArgsUsage: "<repo/name>",
		Action:    prune,
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
				//nolint:gomnd
				Value: 10,
			},
		},
	}
}

func prune(ctx *cli.Context) error {
	client, err := client.New(ctx.StringSlice("server")[0], ctx.String("token"))
	if err != nil {
		return err
	}

	repos, err := getRepos(client)
	if err != nil {
		return err
	}

	buildDuration := ctx.String("older-than")
	buildMin := ctx.Int("keep-min")
	dry := ctx.Bool("dry-run")

	duration, err := time.ParseDuration(buildDuration)
	if err != nil {
		return err
	}

	if dry {
		logrus.Info("dry-run enabled, no data will be removed")
	}

	logrus.Infof("prune builds older than %v, keep min %v", buildDuration, buildMin)

	for _, repo := range repos {
		builds, err := getBuilds(client, repo)
		if err != nil {
			return err
		}

		sort.Slice(builds, func(i, j int) bool {
			return builds[i].Number > builds[j].Number
		})

		if buildCount := len(builds); buildCount > 0 && buildCount > buildMin {
			builds = builds[buildMin:]
			builds = util.Filter(builds, func(b *drone.Build) bool {
				return time.Since(time.Unix(b.Created, 0)) > duration
			})

			logrus.Infof("prune %v/%v builds from '%v'", len(builds), buildCount, repo.Slug)

			if !dry && len(builds) > 0 {
				err := client.BuildPurge(repo.Namespace, repo.Name, int(builds[0].Number+1))
				if err != nil {
					return err
				}
			}
		} else {
			logrus.Infof("skip '%v', number of %v builds lower than min value", repo.Slug, len(builds))
		}
	}

	return nil
}

func getRepos(client drone.Client) ([]*drone.Repo, error) {
	page := 1
	repos := make([]*drone.Repo, 0)

	for {
		repo, err := client.RepoListAll(drone.ListOptions{Page: page, Size: DroneClientPageSize})
		if err != nil {
			return nil, err
		}

		if len(repo) == 0 {
			break
		}

		repos = append(repos, repo...)
		page++
	}

	repos = util.Filter(repos, func(repo *drone.Repo) bool {
		return repo.Active
	})

	return repos, nil
}

func getBuilds(client drone.Client, repo *drone.Repo) ([]*drone.Build, error) {
	page := 1
	builds := make([]*drone.Build, 0)

	for {
		build, err := client.BuildList(repo.Namespace, repo.Name, drone.ListOptions{Page: page, Size: DroneClientPageSize})
		if err != nil {
			return nil, err
		}

		if len(build) == 0 {
			break
		}

		builds = append(builds, build...)
		page++
	}

	return builds, nil
}
