package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

var ListCommand = cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "list diary",
	Action:  ListAction,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "fullpath, f",
			Usage: "list only this year",
		},
		cli.StringFlag{
			Name:  "range, r",
			Usage: "relative date range",
			Value: DefaultDateRange,
		},
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "show all diary's",
		},
	},
}

func ListAction(c *cli.Context) error {
	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}

	diaryList := filterMarkdown(dirWalk(cfg.DiaryDir))
	if !c.Bool("all") {
		filterdList, err := filterDateRange(diaryList, cfg.DiaryDir, c.String("range"))
		if err != nil {
			return err
		}
		diaryList = filterdList
	}

	showPaths := []string{}
	if c.Bool("fullpath") {
		showPaths = diaryList
	} else {
		for _, diaryPath := range diaryList {
			showPaths = append(showPaths, strings.TrimPrefix(diaryPath, cfg.DiaryDir+"/"))
		}
	}

	for _, p := range showPaths {
		fmt.Println(p)
	}
	return nil
}

func filterDateRange(base []string, diaryDirPath string, dateRangeString string) ([]string, error) {
	targetPaths := []string{}
	year, month, day, err := internal.ParseDate(dateRangeString)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	start := now.AddDate(-1*year, -1*month, -1*day)
	dateRange := internal.GetDateRange(start, now)

	targetPaths = []string{}
	for _, d := range dateRange {
		targetPaths = append(targetPaths, internal.DayPath(d, diaryDirPath))
	}

	// Show filtering list
	filteredPaths := []string{}
	for _, p := range base {
		for _, targetPath := range targetPaths {
			if strings.HasPrefix(p, targetPath) {
				filteredPaths = append(filteredPaths, p)
			}
		}
	}

	return filteredPaths, nil
}

func dirWalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirWalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}
