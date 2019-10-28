package cmd

import (
	"fmt"
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

	paths, err := GetDiaryList(cfg.DiaryDir, c.Bool("all"), c.Bool("fullpath"), c.String("range"))
	if err != nil {
		return err
	}

	for _, p := range paths {
		fmt.Println(p)
	}
	return nil
}

func GetDiaryList(diaryDirPath string, isAll bool, isFullPath bool, dateRangeString string) ([]string, error) {
	diaryList := filterMarkdown(internal.Walk(diaryDirPath))
	if !isAll {
		filterdList, err := filterDateRange(diaryList, diaryDirPath, dateRangeString)
		if err != nil {
			return nil, err
		}
		diaryList = filterdList
	}

	showPaths := diaryList
	if !isFullPath {
		showPaths = []string{}
		for _, diaryPath := range diaryList {
			showPaths = append(showPaths, strings.TrimPrefix(diaryPath, diaryDirPath+"/"))
		}
	}
	return showPaths, nil
}

func filterDateRange(base []string, diaryDirPath string, dateRangeString string) ([]string, error) {
	year, month, day, err := internal.ParseDate(dateRangeString)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	start := now.AddDate(-1*year, -1*month, -1*day)
	dateRange := internal.GetDateRange(start, now)

	targetPaths := []string{}
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
