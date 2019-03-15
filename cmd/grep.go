package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

func GrepAction(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return errors.New("The required arguments were not provided: <pattern>")
	}

	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}

	now := time.Now()

	// Relative Range
	files := []string{}
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
			targetPaths = append(targetPaths, internal.DayPath(d, cfg.DiaryDir))
		}

		// Show all list
		all := filterMarkdown(dirWalk(cfg.DiaryDir))

		// Show filtering list
		for _, p := range all {
			for _, targetPath := range targetPaths {
				if strings.HasPrefix(p, targetPath) {
					files = append(files, p)
				}
			}
		}
	} else {
		files = dirWalk(cfg.DiaryDir)
	}

	if len(files) == 0 {
		return errors.New("Not found diary file")
	}

	files = filterMarkdown(files)

	return internal.GrepFiles(cfg.GrepCmd, c.Args().First(), files...)
}

func filterMarkdown(files []string) []string {
	var newfiles []string
	for _, file := range files {
		if strings.HasSuffix(file, ".md") {
			newfiles = append(newfiles, file)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(newfiles)))
	return newfiles
}
