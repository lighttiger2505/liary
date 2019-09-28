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
			Usage: "relative data range",
		},
	},
}

func ListAction(c *cli.Context) error {
	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}
	diaryDirPath := cfg.DiaryDir
	now := time.Now()

	// Relative Range
	if c.String("range") != "" {
		targetPaths := []string{}
		year, month, day, err := internal.ParseDate(c.String("range"))
		if err != nil {
			return err
		}
		start := now.AddDate(-1*year, -1*month, -1*day)
		dateRange := internal.GetDateRange(start, now)

		targetPaths = []string{}
		for _, d := range dateRange {
			targetPaths = append(targetPaths, internal.DayPath(d, diaryDirPath))
		}

		// Show all list
		all := filterMarkdown(dirWalk(diaryDirPath))

		// Show filtering list
		filteredPaths := []string{}
		for _, p := range all {
			for _, targetPath := range targetPaths {
				if strings.HasPrefix(p, targetPath) {
					filteredPaths = append(filteredPaths, p)
				}
			}
		}

		// Remove diary home path
		showPaths := []string{}
		if c.Bool("fullpath") {
			showPaths = filteredPaths
		} else {
			for _, diaryPath := range filteredPaths {
				showPaths = append(showPaths, strings.TrimPrefix(diaryPath, cfg.DiaryDir+"/"))
			}
		}

		// Show filtering list
		for _, p := range showPaths {
			fmt.Println(p)
		}
		return nil
	}

	// Show all list
	all := filterMarkdown(dirWalk(diaryDirPath))

	// Remove diary home path
	showPaths := []string{}
	if c.Bool("fullpath") {
		showPaths = all
	} else {
		for _, diaryPath := range all {
			showPaths = append(showPaths, strings.TrimPrefix(diaryPath, cfg.DiaryDir+"/"))
		}
	}

	for _, diaryPath := range showPaths {
		fmt.Println(diaryPath)
	}
	return nil
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
