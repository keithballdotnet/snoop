package main

import (
	"fmt"
	"net/http"

	gitlab "github.com/xanzy/go-gitlab"
)

func getMergeRequest(pid, mrid int) error {
	//156, 55
	mr, resp, err := git.MergeRequests.GetMergeRequest(pid, mrid)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: got status code: %v", resp.StatusCode)
	}

	project, err := getProject(mr.ProjectID)
	if err != nil {
		return fmt.Errorf("error: unable to find project: %v", mr.ProjectID)
	}
	fmt.Printf("MR: %s - %s - %v\n", mr.Author.Name, mr.Title, project.Name)

	return nil
}

func getMergeRequestCommits() error {

	allCommits := []*gitlab.Commit{}

	opts := &gitlab.GetMergeRequestCommitsOptions{PerPage: 100, Page: 1}

	commits, resp, err := git.MergeRequests.GetMergeRequestCommits(154, 32, opts)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: got status code: %v", resp.StatusCode)
	}

	for _, c := range commits {
		allCommits = append(allCommits, c)
	}

	for page := 2; page <= resp.TotalPages; page++ {
		opts.Page = page
		commits, _, err := git.MergeRequests.GetMergeRequestCommits(154, 32, opts)
		if err != nil {
			return err
		}

		// if resp.StatusCode != http.StatusOK {
		// 	return fmt.Errorf("error: got status code: %v", resp.StatusCode)
		// }

		for _, c := range commits {
			allCommits = append(allCommits, c)
		}
	}

	// for _, c := range allCommits {
	// 	fmt.Printf("Commit: %s - %v\n", c.CommitterName, c.CommittedDate.Format("2006-01-02"))
	// }

	// fmt.Printf("Commits: %v\n", len(allCommits))

	generateReport(allCommits)

	return nil
}

type workDay struct {
	Date    string
	Commits map[string]int
}

func generateReport(commits []*gitlab.Commit) {

	// Reverse slice...
	for i := len(commits)/2 - 1; i >= 0; i-- {
		opp := len(commits) - 1 - i
		commits[i], commits[opp] = commits[opp], commits[i]
	}

	workDays := []workDay{}

	// last commit day
	lastCommitDate := commits[0].CommittedDate.Format("2006-01-02")
	commitCountMap := make(map[string]int)
	day := workDay{Date: lastCommitDate, Commits: commitCountMap}
	for i := 0; i < len(commits); i++ {
		currentCommitDate := commits[i].CommittedDate.Format("2006-01-02")
		// Date change...
		if lastCommitDate != currentCommitDate {
			//fmt.Printf("DateChange: Last - %s Current - %s\n", lastCommitDate, currentCommitDate)
			workDays = append(workDays, day)
			lastCommitDate = currentCommitDate
			commitCountMap = make(map[string]int)
			day = workDay{Date: lastCommitDate, Commits: commitCountMap}

		}
		// Increase commit count for user on this day
		count, ok := day.Commits[commits[i].CommitterEmail]
		if ok {
			day.Commits[commits[i].CommitterEmail] = count + 1
		} else {
			day.Commits[commits[i].CommitterEmail] = 1
		}

	}

	// Append the last work day...
	workDays = append(workDays, day)

	fte := 0
	for _, d := range workDays {
		fmt.Printf("Day: %s\n", d.Date)
		for key, value := range d.Commits {
			fmt.Printf("Commits: %v Commiter: %s\n", value, key)
			fte = fte + 1
		}
	}
	fmt.Printf("Total Work Days: %v\n", len(workDays))
	fmt.Printf("Total FTE Days: %v\n", fte)

	// firstCommitTime := commits[0].CommittedDate
	// lastCommitTime := commits[len(commits)-1].CommittedDate
	// calendarDays := lastCommitTime.Sub(*firstCommitTime).Hours()
	// fmt.Printf("Total Calendar Days: %v\n", calendarDays)

}

func getMergeRequests(state string) error {

	//fmt.Println("getting all merge requests. this may take a while...")

	mrs, resp, err := git.MergeRequests.ListMergeRequests(
		&gitlab.ListMergeRequestsOptions{
			Scope: gitlab.String("all"),
			State: gitlab.String(state),
			ListOptions: gitlab.ListOptions{
				PerPage: 10000,
				Page:    1,
			}})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: got status code: %v", resp.StatusCode)
	}

	for _, mr := range mrs {
		// Show only in review MRs
		if mr.WorkInProgress {
			continue
		}
		project, err := getProject(mr.ProjectID)
		if err != nil {
			return fmt.Errorf("error: unable to find project: %v", mr.ProjectID)
		}
		fmt.Printf("MR: %s - %s - %v\n", mr.Author.Name, mr.Title, project.Name)
	}

	return nil
}