package cmd

import (
	"errors"
	"fmt"

	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

func ConfigAction(c *cli.Context) error {
	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}

	if c.String("get") != "" {
		switch c.String("get") {
		case "diarydir":
			fmt.Println(cfg.DiaryDir)
		case "editor":
			fmt.Println(cfg.Editor)
		case "grepcmd":
			fmt.Println(cfg.GrepCmd)
		default:
			return errors.New("key does not contain a section")
		}
		return nil
	}
	if c.String("get-all") != "" {
	}

	internal.OpenEditor(cfg.Editor, cfg.Path())
	return nil
}
