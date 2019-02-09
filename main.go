package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lighttiger2505/liary/internal"
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
			Action:  edit,
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
			Action:  appendRun,
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
			Action:  list,
		},
	}
	// 	app.Action = run
	return app
}

func edit(c *cli.Context) error {
	// Getting time for target diary
	date := c.String("date")
	before := c.Int("before")
	after := c.Int("after")
	targetTime, err := getTargetTime(date, before, after)
	if err != nil {
		return err
	}

	// Getting diary path
	suffix := suffixJoin(c.String("suffix"))
	targetPath, err := internal.DiaryPath(targetTime, internal.DiaryDirPath(), suffix)
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
	return Open(targetPath)
}

func list(c *cli.Context) error {
	// Show diary file list
	return ListAll()
}

func appendRun(c *cli.Context) error {
	// Getting time for target diary
	date := c.String("date")
	before := c.Int("before")
	after := c.Int("after")
	targetTime, err := getTargetTime(date, before, after)
	if err != nil {
		return err
	}

	// Getting diary path
	suffix := suffixJoin(c.String("suffix"))
	targetPath, err := internal.DiaryPath(targetTime, internal.DiaryDirPath(), suffix)
	if err != nil {
		return err
	}

	// Make directory
	targetDirPath := filepath.Dir(targetPath)
	if err := internal.MakeDir(targetDirPath); err != nil {
		return err
	}

	appendVal, err := internal.GetAppendValue(c.Args())
	if err != nil {
		return fmt.Errorf("Failed get append value %s", err)
	}

	code := c.Bool("code")
	lang := c.String("language")
	numLineBefore := c.Int("before-append")
	numLineAfter := c.Int("after-append")
	if appendVal != "" {
		if code {
			return AppendCodeBlock(targetPath, appendVal, numLineBefore, numLineAfter, lang)
		}
		return Append(targetPath, appendVal, numLineBefore, numLineAfter)
	}
	return nil
}

func suffixJoin(val string) string {
	words := strings.Fields(val)
	return strings.Join(words, "_")
}

func getTargetTime(date string, before, after int) (time.Time, error) {
	if date != "" {
		return time.Parse("2006-01-02", date)
	}
	now := time.Now()
	return internal.UpDonwDate(now, before, after)
}
