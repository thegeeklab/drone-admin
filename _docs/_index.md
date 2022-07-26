---
title: drone-admin
---

[![Build Status](https://img.shields.io/drone/build/thegeeklab/drone-admin?logo=drone&server=https%3A%2F%2Fdrone.thegeeklab.de)](https://drone.thegeeklab.de/thegeeklab/drone-admin)
[![Docker Hub](https://img.shields.io/badge/dockerhub-latest-blue.svg?logo=docker&logoColor=white)](https://hub.docker.com/r/thegeeklab/drone-admin)
[![Quay.io](https://img.shields.io/badge/quay-latest-blue.svg?logo=docker&logoColor=white)](https://quay.io/repository/thegeeklab/drone-admin)
[![Go Report Card](https://goreportcard.com/badge/github.com/thegeeklab/drone-admin)](https://goreportcard.com/report/github.com/thegeeklab/drone-admin)
[![GitHub contributors](https://img.shields.io/github/contributors/thegeeklab/drone-admin)](https://github.com/thegeeklab/drone-admin/graphs/contributors)
[![Source: GitHub](https://img.shields.io/badge/source-github-blue.svg?logo=github&logoColor=white)](https://github.com/thegeeklab/drone-admin)
[![License: Apache-2.0](https://img.shields.io/github/license/thegeeklab/drone-admin)](https://github.com/thegeeklab/drone-admin/blob/main/LICENSE)

Admin Tools for [Drone](https://github.com/drone/drone).

<!-- prettier-ignore-start -->
<!-- spellchecker-disable -->
{{< toc >}}
<!-- spellchecker-enable -->
<!-- prettier-ignore-end -->

## Build

Build the binary with the following command:

```Shell
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

make build
```

Build the Docker image with the following command:

```Shell
docker build --file docker/Dockerfile.amd64 --tag thegeeklab/drone-admin .
```

## Usage

```Shell
drone-admin --help
NAME:
   drone-admin - drone admin tools

USAGE:
   drone-admin [global options] command [command options] [arguments...]

VERSION:
   00d2c63

COMMANDS:
   build       manage build
   autoscaler  manage autoscaler
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dry-run                 disable api calls (default: false) [$DRONE_ADMIN_DRY_RUN]
   --help, -h                show help (default: false)
   --log-level value         log level (default: "info") [$DRONE_ADMIN_LOG_LEVEL]
   --server value, -s value  server address  (accepts multiple inputs) [$DRONE_ADMIN_SERVER]
   --token value, -t value   server auth token [$DRONE_ADMIN_TOKEN]
   --version, -v             print the version (default: false)
```

### Prune Builds

The build prune subcommand removes all builds older than the specified time limit (seconds, minutes or hours) from all repositories. This command should be used with an admin token, otherwise only repositories to which the user has access can be removed. It is also possible to define a minimum number of builds for each repository, even if the time limit has been exceeded. This can be useful to keep a minimum number of builds for repositories with low frequency.

```Shell
drone-admin --token my-secret-token --server https://drone.excample.com build prune --older-than 720h
INFO[0001] prune builds older than 720h, keep min 10
INFO[0001] skip 'example/repo_1', number of 9 builds lower than min value
INFO[0002] prune 1/105 builds from 'example/demo'
INFO[0002] prune 0/56 builds from 'example/cool_project'
```

### Cleanup autoscaler agents

When using the autoscaler, agents sometimes remain in error state in the DB (even if the Drone CI Reaper is enabled). This command tries the destroy agents in error state two times and forces it on the third attempt. For this command the `--server` flag must be set to the address of the autoscaler server(s).

```Shell
drone-admin --token my-secret-token --server https://drone-scaler.excample.com autoscaler reaper
INFO[0000] lookup agents in error state                  error=1 ok=1 server="https://drone-scaler.excample.com"
INFO[0000] destroy agent                                 agent=agent-G8hHyA0A force=false server="https://drone-scaler.excample.com" triage=1
```
