package cmd

import (
	"fmt"

	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

func ConfigAction(c *cli.Context) error {
	fmt.Println("getconfig start")
	cfg, err := internal.GetConfig()
	fmt.Println("getconfig end")
	if err != nil {
		return err
	}
	internal.OpenEditor(cfg.Editor, cfg.Path())
	return nil
}
