package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/snabb/isoweek"
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
	}

	app.Commands = []cli.Command{
		{
			Name: "merge",
			Flags: append(defaultFlags,
				cli.StringFlag{
					Name:  "merge_id, m",
					Usage: "merge request id",
				}),
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
			Name: "merge_commits",
			Flags: append(defaultFlags,
				cli.StringFlag{
					Name:  "merge_id, m",
					Usage: "merge request id",
				}),
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
			Name: "project_merges",
			Flags: append(append(defaultFlags,
				cli.IntFlag{
					Name:  "weeks, w",
					Usage: "# of weeks to go back",
				}),
				cli.StringFlag{
					Name:  "branch, b",
					Usage: "target branch of MRs",
				}),
			Usage: "get project merge requests",
			Action: func(c *cli.Context) error {
				exitErr := setGit()
				if exitErr != nil {
					return exitErr
				}

				var updateAfter *time.Time
				if c.IsSet("weeks") {
					dur := time.Duration(time.Duration(c.Int("weeks")) * 7 * 24 * time.Hour)
					t := time.Now().Add(-dur)
					// Go to start of that ISO week...
					yr, wk := t.ISOWeek()
					t = isoweek.StartTime(yr, wk, time.UTC)
					updateAfter = &t
					//fmt.Printf("Going to use: %s\n", updateAfter.String())
				}

				branch := new(string)
				if c.IsSet("branch") {
					*branch = c.String("branch")
				}

				err := getProjectMergeRequests(c.Int("project_id"), branch, updateAfter)
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
