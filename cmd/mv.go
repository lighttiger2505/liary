package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

func MoveAction(c *cli.Context) error {
	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}

	args := c.Args()
	realSource := filepath.Join(cfg.DiaryDir, args[0])
	realDist := filepath.Join(cfg.DiaryDir, args[1])

	if !isFileExist(realSource) {
		return fmt.Errorf("missing source file operand after, '%v'", realSource)
	}

	if err := os.Rename(realSource, realDist); err != nil {
		return fmt.Errorf("failed move file, '%v'", err.Error())
	}

	return nil
}

func isFileExist(fPath string) bool {
	_, err := os.Stat(fPath)
	return err == nil || !os.IsNotExist(err)
}
