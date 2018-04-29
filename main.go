package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
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
	targetTime, err := targetTime(date, before, after)
	if err != nil {
		return err
	}

	// Getting diary path
	file := c.String("file")
	suffix := c.String("suffix")
	targetPath := ""
	if file != "" {
		targetPath = file
	} else {
		targetPath, err = diaryPath(targetTime, diaryDirPath(), suffix)
		if err != nil {
			return err
		}
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

	appendVal, err := getAppendValue(c.Args())
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

func diaryDirPath() string {
	home, _ := homedir.Dir()
	diaryDirPath := filepath.Join(home, "diary")
	return diaryDirPath
}

func MonthPath(targetTime time.Time, dirPath string) string {
	year, month, _ := targetTime.Date()
	result := filepath.Join(
		dirPath,
		fmt.Sprintf("%02d", year),
		fmt.Sprintf("%02d", int(month)),
	)
	return result
}

func DayPath(targetTime time.Time, dirPath string) string {
	year, month, day := targetTime.Date()
	diaryPath := filepath.Join(
		dirPath,
		fmt.Sprintf("%02d", year),
		fmt.Sprintf("%02d", int(month)),
		fmt.Sprintf("%02d", day),
	)
	return diaryPath
}

func diaryPath(targetTime time.Time, dirPath string, suffix string) (string, error) {
	year, month, day := targetTime.Date()

	var filename string
	if suffix != "" {
		filename = fmt.Sprintf("%s-%s.md", fmt.Sprintf("%02d", day), suffix)
	} else {
		filename = fmt.Sprintf("%s.md", fmt.Sprintf("%02d", day))
	}

	diaryPath := filepath.Join(
		dirPath,
		fmt.Sprintf("%02d", year),
		fmt.Sprintf("%02d", int(month)),
		filename,
	)
	return diaryPath, nil
}

func targetTime(date string, before, after int) (time.Time, error) {
	now := time.Now()
	if date != "" {
		now, err := time.Parse("2006-01-02", date)
		if err != nil {
			return now, err
		}
	}
	if before != 0 {
		return now.AddDate(0, 0, -1*before), nil
	}
	if after != 0 {
		return now.AddDate(0, 0, after), nil
	}
	return now, nil
}

func isFileExist(fPath string) bool {
	_, err := os.Stat(fPath)
	return err == nil || !os.IsNotExist(err)
}

func searchBeforeDate(startDate time.Time, dirPath string) (time.Time, error) {

	for i := 0; i < 30; i++ {
		date := startDate.AddDate(0, 0, -1*i)
		path, err := diaryPath(date, diaryDirPath(), "")
		if err != nil {
			return time.Time{}, err
		}

		if isFileExist(path) {
			return date, nil
		}
	}
	return time.Time{}, fmt.Errorf("Not found before diary")
}

func getAppendValue(args []string) (string, error) {
	var val string
	if terminal.IsTerminal(0) {
		if len(args) > 0 {
			val = args[0]
		} else {
			val = ""
		}
	} else {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("Failed make diary file. %s", err.Error())
		}
		val = string(b)
	}
	return val, nil
}
