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
		cmd.EditCommand,
		cmd.AppendCommand,
		cmd.ListCommand,
		cmd.GrepCommand,
		cmd.ConfigCommand,
		cmd.MoveCommand,
	}
	return app
}
