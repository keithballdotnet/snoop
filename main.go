package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
	gitlab "github.com/xanzy/go-gitlab"
)

var git *gitlab.Client

func main() {
	app := cli.NewApp()
	app.Name = "snoop"
	app.Usage = "gather some metrics from gitlab"
	app.Version = "1.0.0"

	var server, token string
	setGit := func() *cli.ExitError {
		if token == "" {
			return cli.NewExitError("no token specified", 1)
		}
		if server == "" {
			return cli.NewExitError("no server specified", 1)
		}

		git = gitlab.NewClient(nil, token)
		host := fmt.Sprintf("https://%s/api/v4", server)
		git.SetBaseURL(host)

		return nil
	}

	defaultFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "server, s",
			Usage:       "gitlab server url",
			Destination: &server,
			EnvVar:      "SNOOP_SERVER",
		},
		cli.StringFlag{
			Name:        "token, t",
			Usage:       "auth token",
			Destination: &token,
			EnvVar:      "SNOOP_TOKEN",
		},
		cli.IntFlag{
			Name:  "project_id, p",
			Usage: "project id",
		},
		cli.StringFlag{
			Name:  "merge_id, m",
			Usage: "merge request id",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "merge",
			Flags: defaultFlags,
			Usage: "get merge info",
			Action: func(c *cli.Context) error {

				exitErr := setGit()
				if exitErr != nil {
					return exitErr
				}

				//fmt.Printf("Merge request: %v %v\n", c.Int("project_id"), c.Int("merge_id"))

				err := getMergeRequest(c.Int("project_id"), c.Int("merge_id"))
				if err != nil {
					return cli.NewExitError("error: "+err.Error(), 1)
				}

				return nil
			},
		},
		{
			Name:  "merge_commits",
			Flags: defaultFlags,
			Usage: "get merge commits",
			Action: func(c *cli.Context) error {

				exitErr := setGit()
				if exitErr != nil {
					return exitErr
				}

				err := getMergeRequestCommits(c.Int("project_id"), c.Int("merge_id"))
				if err != nil {
					return cli.NewExitError("error: "+err.Error(), 1)
				}

				return nil
			},
		},
		{
			Name:  "project_merges",
			Flags: defaultFlags,
			Usage: "get project merge requests",
			Action: func(c *cli.Context) error {
				exitErr := setGit()
				if exitErr != nil {
					return exitErr
				}

				err := getProjectMergeRequest(c.Int("project_id"))
				if err != nil {
					return cli.NewExitError("error: "+err.Error(), 1)
				}
				return nil
			},
		},
	}

	// // Verify input
	// app.Action = func(c *cli.Context) error {
	// 		fmt.Printf("server: %s\n", server)

	// 	return nil
	// }

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
