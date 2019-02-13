package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

func EditAction(c *cli.Context) error {
	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}

	// Getting time for target diary
	date := c.String("date")
	before := c.Int("before")
	after := c.Int("after")
	targetTime, err := internal.GetTargetTime(date, before, after)
	if err != nil {
		return err
	}

	// Getting diary path
	suffix := ""
	args := c.Args()
	if len(args) > 0 {
		suffix = suffixJoin(args[0])
	}
	targetPath, err := internal.DiaryPath(targetTime, cfg.DiaryDir, suffix)
	if err != nil {
		return err
	}

	// Show diary file path
	path := c.Bool("path")
	if path {
		fmt.Println(targetPath)
		return nil
	}

	// Make directory
	targetDirPath := filepath.Dir(targetPath)
	if err := internal.MakeDir(targetDirPath); err != nil {
		return err
	}

	// Open text editor
	return internal.OpenEditor(cfg.Editor, targetPath)
}

func suffixJoin(val string) string {
	words := strings.Fields(val)
	return strings.Join(words, "_")
}
