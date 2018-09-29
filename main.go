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
	app.Flags = []cli.Flag{
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
			EnvVar:      "SNOOP_AUTH_TOKEN",
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "merge",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "project_id, p",
					Usage: "project id",
				},
				cli.StringFlag{
					Name:  "merge_id, m",
					Usage: "merge request id",
				},
			},
			Usage: "get merge info",
			Action: func(c *cli.Context) error {
				fmt.Printf("Merge request: %v %v\n", c.Int("project_id"), c.Int("merge_id"))

				err := getMergeRequest(c.Int("project_id"), c.Int("merge_id"))
				if err != nil {
					return cli.NewExitError("error: "+err.Error(), 1)
				}

				return nil
			},
		},
	}

	// Verify input
	app.Action = func(c *cli.Context) error {
		if token == "" {
			return cli.NewExitError("no token specified", 1)
		}
		if server == "" {
			return cli.NewExitError("no server specified", 1)
		}
		return nil
	}

	git = gitlab.NewClient(nil, token)
	host := fmt.Sprintf("https://%s/api/v4", server)
	git.SetBaseURL(host)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
