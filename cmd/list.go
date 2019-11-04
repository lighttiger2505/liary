package cmd

import (
	"fmt"

	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

var ListCommand = cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "list diary",
	Action:  ListAction,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "fullpath, f",
			Usage: "list only this year",
		},
		cli.StringFlag{
			Name:  "range, r",
			Usage: "relative date range",
			Value: DefaultDateRange,
		},
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "show all diary's",
		},
	},
}

func ListAction(c *cli.Context) error {
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

	paths, err := internal.GetDiaryList(workspace, c.Bool("all"), c.Bool("fullpath"), c.String("range"))
	if err != nil {
		return err
	}

	for _, p := range paths {
		fmt.Println(p)
	}
	return nil
}
