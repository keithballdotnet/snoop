package main

import (
	"errors"
	"fmt"
	"net/http"

	gitlab "github.com/xanzy/go-gitlab"
)

var cachedProjects = make(map[int]*gitlab.Project)

func getProjects() error {

	//fmt.Println("getting all merge requests. this may take a while...")

	opt := &gitlab.ListProjectsOptions{
		OrderBy:    gitlab.String("id"),
		Simple:     gitlab.Bool(false),
		Membership: gitlab.Bool(false),
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
	}

	projects, resp, err := git.Projects.ListProjects(opt)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: got status code: %v", resp.StatusCode)
	}

	// load project cache
	for _, project := range projects {
		//fmt.Printf("Project: %s %v\n", project.Name, project.ID)
		cachedProjects[project.ID] = project
	}

	for page := 2; page < resp.TotalPages; page++ {
		opt.ListOptions.Page = page
		projects, _, err := git.Projects.ListProjects(opt)
		if err != nil {
			return err
		}
		// if resp.StatusCode != http.StatusOK {
		// 	return fmt.Errorf("error: got status code: %v", resp.StatusCode)
		// }

		// load project cache
		for _, project := range projects {
			//fmt.Printf("Project: %s %v\n", project.Name, project.ID)
			cachedProjects[project.ID] = project
		}
	}

	return nil
}

// getProject will get a project
func getProject(id int) (*gitlab.Project, error) {
	p, ok := cachedProjects[id]
	if !ok {
		// No in cache
		project, resp, err := git.Projects.GetProject(id)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error: got status code: %v", resp.StatusCode)
		}
		if project == nil {
			return nil, errors.New("unable to find project")
		}

		return project, nil
	}

	return p, nil
}
