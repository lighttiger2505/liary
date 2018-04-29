package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

const (
	ExitCodeOK        int = iota //0
	ExitCodeFileError int = iota //2
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
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "suffix, x",
			Usage: "Diary file suffix",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Specified file",
		},
		cli.BoolFlag{
			Name:  "list, l",
			Usage: "Show diary file list",
		},
		cli.BoolFlag{
			Name:  "today, t",
			Usage: "Show diary file list on today",
		},
		cli.BoolFlag{
			Name:  "path, p",
			Usage: "Show diary file path",
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
	}
	app.Action = run
	return app
}

func run(c *cli.Context) error {
	// Getting time for target diary
	date := c.String("date")
	before := c.Int("before")
	after := c.Int("after")
	targetTime, err := TargetTime(date, before, after)
	if err != nil {
		return err
	}

	// Getting diary path
	suffix := c.String("suffix")
	targetPath, err := DiaryPath(targetTime, diaryDirPath(), suffix)
	if err != nil {
		return err
	}

	// Show diary file list
	list := c.Bool("list")
	if list {
		if c.GlobalBool("date") || c.GlobalBool("before") || c.GlobalBool("after") || c.GlobalBool("today") {
			return ListTargetDate(targetTime)
		}
		return ListAll()
	}

	// Show diary file path
	path := c.Bool("path")
	if path {
		fmt.Println(targetPath)
		return nil
	}

	// Make directory
	targetDirPath := filepath.Dir(targetPath)
	if !isFileExist(targetDirPath) {
		if err := os.MkdirAll(targetDirPath, 0755); err != nil {
			return fmt.Errorf("Failed make diary dir. %s", err.Error())
		}
	}

	appendVal, err := GetAppendValue(c.Args())
	if err != nil {
		return fmt.Errorf("Failed get append value %s", err)
	}

	if appendVal != "" {
		// Append content
		return Append(targetPath, appendVal)
	}

	// Open text editor
	return Open(targetPath)
}
