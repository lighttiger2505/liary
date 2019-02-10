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

func ListAction(c *cli.Context) error {
	now := time.Now()
	diaryDirPath := internal.DiaryDirPath()

	// Absolute Range
	targetPaths := []string{}
	if c.Bool("today") {
		targetPaths = append(targetPaths, internal.DayPath(now, diaryDirPath))
	} else if c.Bool("week") {
		for _, d := range internal.GetWeakDays(now) {
			targetPaths = append(targetPaths, internal.DayPath(d, diaryDirPath))
		}
	} else if c.Bool("month") {
		targetPaths = append(targetPaths, internal.MonthPath(now, diaryDirPath))
	} else if c.Bool("year") {
		targetPaths = append(targetPaths, internal.YearPath(now, diaryDirPath))
	}

	// Relative Range
	if c.String("range") != "" {
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
	}

	diaryPaths := dirWalk(diaryDirPath)

	// Show all list
	if len(targetPaths) == 0 {
		for _, diaryPath := range diaryPaths {
			fmt.Println(diaryPath)
		}
		return nil
	}

	// Show filtering list
	for _, diaryPath := range diaryPaths {
		for _, targetPath := range targetPaths {
			if strings.HasPrefix(diaryPath, targetPath) {
				fmt.Println(diaryPath)
			}
		}
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
