package cmd

import (
	"errors"
	"strings"
	"time"

	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

var GrepCommand = cli.Command{
	Name:    "grep",
	Aliases: []string{"g"},
	Usage:   "grep diary",
	Action:  GrepAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "range, r",
			Usage: "relative date range",
			Value: DefaultDateRange,
		},
	},
}

func GrepAction(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return errors.New("The required arguments were not provided: <pattern>")
	}

	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}

	workspace := cfg.DiaryDir
	workspaceFlag := c.GlobalString("workspace")
	if workspaceFlag != "" {
		w, err := cfg.GetWorkSpace(workspaceFlag)
		if err != nil {
			return err
		}
		workspace = w
	}

	// Relative Range
	files := []string{}
	if c.String("range") != "" {
		year, month, day, err := internal.ParseDate(c.String("range"))
		if err != nil {
			return err
		}
		now := time.Now()
		start := now.AddDate(-1*year, -1*month, -1*day)
		dateRange := internal.GetDateRange(start, now)

		targetPaths := []string{}
		for _, d := range dateRange {
			targetPaths = append(targetPaths, internal.DayPath(d, workspace))
		}

		// Show all list
		all := internal.FilterMarkdown(internal.Walk(workspace))

		// Show filtering list
		for _, p := range all {
			for _, targetPath := range targetPaths {
				if strings.HasPrefix(p, targetPath) {
					files = append(files, p)
				}
			}
		}
	} else {
		files = internal.Walk(workspace)
	}

	if len(files) == 0 {
		return errors.New("Not found diary file")
	}

	files = internal.FilterMarkdown(files)

	return internal.GrepFiles(cfg.GrepCmd, c.Args().First(), files...)
}
