package cmd

import (
	"errors"
	"sort"
	"strings"

	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

func GrepAction(c *cli.Context) error {
	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}

	if len(c.Args()) == 0 {
		return errors.New("The required arguments were not provided: <pattern>")
	}

	files := dirWalk(cfg.DiaryDir)
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
