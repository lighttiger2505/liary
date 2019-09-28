package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lighttiger2505/liary/internal"
	"github.com/lighttiger2505/liary/internal/ui"
	"github.com/urfave/cli"
)

var MoveCommand = cli.Command{
	Name:      "mv",
	Aliases:   []string{"m"},
	Usage:     "move diary",
	UsageText: "liary mv [command options...] <source file> <dist file>",
	Action:    MoveAction,
	Flags:     []cli.Flag{},
}

func MoveAction(c *cli.Context) error {
	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}
	args := c.Args()

	// Fetching source file path
	// Check both absolute and relative paths
	source := args[0]
	var absSourcePath string
	if !filepath.IsAbs(source) {
		absSourcePath = filepath.Join(cfg.DiaryDir, source)
	}
	if !isFileExist(absSourcePath) {
		return fmt.Errorf("missing source file operand after, '%v'", absSourcePath)
	}

	// Fetching dist file path
	// Check both absolute and relative paths
	dist := args[1]
	var absDistPath string
	if !filepath.IsAbs(dist) {
		absDistPath = filepath.Join(cfg.DiaryDir, dist)
	}
	if isFileExist(absDistPath) {
		in, err := ui.Ask(fmt.Sprintf("are you sure you want to overwrite this file, '%s'? (y/n)", absDistPath))
		if err != nil {
			return err
		}
		if !ui.CheckAnswerYes(in) {
			return nil
		}
	}

	if err := os.Rename(absSourcePath, absDistPath); err != nil {
		return fmt.Errorf("failed move file, '%v'", err.Error())
	}

	return nil
}

func isFileExist(fPath string) bool {
	_, err := os.Stat(fPath)
	return err == nil || !os.IsNotExist(err)
}
