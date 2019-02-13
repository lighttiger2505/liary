package cmd

import (
	"errors"

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

	return internal.GrepFiles(cfg.GrepCmd, c.Args().First(), cfg.DiaryDir)
}
