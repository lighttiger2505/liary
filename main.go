package main

import (
	"fmt"
	"os"

	"github.com/lighttiger2505/liary/cmd"
	"github.com/urfave/cli"
)

const (
	ExitCodeOK    int = iota //0
	ExitCodeError int = iota //1
)

func main() {
	err := newApp().Run(os.Args)
	var exitCode = ExitCodeOK
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		exitCode = ExitCodeError
	}
	os.Exit(exitCode)
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "liary"
	app.HelpName = "liary"
	app.Usage = "liary is fastest cli tool for create a diary."
	app.UsageText = "liary [options] [write content for diary]"
	app.Version = "0.0.1"
	app.Author = "lighttiger2505"
	app.Email = "lighttiger2505@gmail.com"
	app.Commands = []cli.Command{
		cli.Command{
			Name:    "edit",
			Aliases: []string{"e"},
			Usage:   "edit diary",
			Action:  cmd.EditAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "suffix, x",
					Usage: "Diary file suffix",
				},
				cli.StringFlag{
					Name:  "date, d",
					Usage: "Specified date",
				},
				cli.IntFlag{
					Name:  "before, b",
					Usage: "Specified before diary by day",
				},
				cli.IntFlag{
					Name:  "after, a",
					Usage: "Specified after diary by day",
				},
			},
		},
		cli.Command{
			Name:    "append",
			Aliases: []string{"a"},
			Usage:   "grep diary",
			Action:  cmd.AppendAction,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "code, c",
					Usage: "append code block",
				},
				cli.StringFlag{
					Name:  "language, g",
					Usage: "code block language",
				},
				cli.IntFlag{
					Name:  "before-append, B",
					Usage: "NUM of blank line to add before content to be append",
					Value: 1,
				},
				cli.IntFlag{
					Name:  "after-append, A",
					Usage: "NUM of blank line to add after content to be append",
					Value: 1,
				},
			},
		},
		cli.Command{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list diary",
			Action:  cmd.ListAction,
		},
	}
	// 	app.Action = run
	return app
}
