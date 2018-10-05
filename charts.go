package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/wcharczuk/go-chart"
	gitlab "github.com/xanzy/go-gitlab"
)

func getProjectMergeRequestsChart(project *gitlab.Project, db map[string]weeklyMergeRequestInfo) error {

	// Sort the keys, so output is nice
	var keys []string
	for k := range db {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	values := []chart.Value{}

	maxValue := 1
	for _, key := range keys {
		entry := db[key]
		total := /*entry.CountOfOpened + entry.CountOfClosed +*/ entry.CountOfMerged
		values = append(values, chart.Value{Value: float64(total), Label: key})
		if total > maxValue {
			maxValue = total
		}
	}

	chartWidth := (50 * len(db)) + 150

	sbc := chart.BarChart{
		Title:      project.Name + " - Merged MRs",
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top:    40,
				Bottom: 40,
			},
		},
		Height:   512,
		Width:    chartWidth,
		BarWidth: 50,
		XAxis: chart.Style{
			Show: true,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
			Range: &chart.ContinuousRange{
				Min: 0,
				Max: float64(maxValue),
			},
		},
		Bars: values,
	}

	tmpfile, err := ioutil.TempFile("", "snoop")
	if err != nil {
		return err
	}

	err = sbc.Render(chart.PNG, tmpfile)
	if err != nil {
		return err
	}

	err = tmpfile.Close()
	if err != nil {
		return err
	}

	pngFile := tmpfile.Name() + ".png"
	err = os.Rename(tmpfile.Name(), pngFile)
	if err != nil {
		return err
	}

	fmt.Printf("Wrote chart here: %s\n", pngFile)

	return nil
}
