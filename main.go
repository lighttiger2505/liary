package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

const (
	ExitCodeOK        int = iota //0
	ExitCodeError     int = iota //1
	ExitCodeFileError int = iota //2
)

func main() {
	err := newApp().Run(os.Args)
	var exitCode = ExitCodeOK
	if err != nil {
		fmt.Println(os.Stderr, err.Error())
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
			Name:  "file, f",
			Usage: "Specified file",
		},
		cli.BoolFlag{
			Name:  "list, l",
			Usage: "Show diary file list",
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

	// Show diary file list
	list := c.Bool("list")
	if list {
		diaryDirPath := diaryDirPath()
		diaryPaths := dirWalk(diaryDirPath)
		for _, diaryPath := range diaryPaths {
			fmt.Println(diaryPath)
		}
		return nil
	}

	// Getting time for target diary
	before := c.Int("before")
	after := c.Int("after")
	targetTime := targetTime(before, after)

	// Getting diary path
	year, month, day := targetTime.Date()
	diaryPath, err := diaryPath(strconv.Itoa(year), strconv.Itoa(int(month)), strconv.Itoa(day))

	// Show diary file path
	path := c.Bool("path")
	if path {
		if err != nil {
			return err
		}
		fmt.Println(diaryPath)
		return nil
	}

	return nil
}

func diaryDirPath() string {
	home, _ := homedir.Dir()
	diaryDirPath := filepath.Join(home, "diary")
	return diaryDirPath
}

func diaryPath(year, month, day string) (string, error) {
	diaryDirPath := diaryDirPath()
	diaryPath := filepath.Join(diaryDirPath, year, month, fmt.Sprintf("%s.md", day))
	return diaryPath, nil
}

func targetTime(before, after int) time.Time {
	now := time.Now()
	if before != 0 {
		return now.AddDate(0, 0, before)
	}
	if after != 0 {
		return now.AddDate(0, 0, -1*after)
	}
	return now
}

func dirWalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirWalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}
