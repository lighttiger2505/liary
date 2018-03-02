package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
	date := c.String("date")
	before := c.Int("before")
	after := c.Int("after")
	targetTime, err := targetTime(date, before, after)
	if err != nil {
		return err
	}

	// Getting diary path
	file := c.String("file")
	targetPath := ""
	if file != "" {
		targetPath = file
	} else {
		targetPath, err = diaryPath(targetTime, diaryDirPath())
		if err != nil {
			return err
		}
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

	// Make diary file
	if !isFileExist(targetPath) {
		if err := makeFile(targetPath); err != nil {
			return fmt.Errorf("Failed make diary file. %s", err.Error())
		}
	}

	args := c.Args()
	if len(args) > 0 {
		// Append content
		file, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("Failed append diary. %s", err.Error())
		}
		defer file.Close()
		fmt.Fprintln(file, args[0])
	} else {
		// Open text editor
		editorEnv := os.Getenv("EDITOR")
		if editorEnv == "" {
			editorEnv = "vim"
		}
		err = openEditor(editorEnv, targetPath)
		if err != nil {
			return fmt.Errorf("Failed open editor. %s", err.Error())
		}
	}

	return nil
}

func diaryDirPath() string {
	home, _ := homedir.Dir()
	diaryDirPath := filepath.Join(home, "diary")
	return diaryDirPath
}

func diaryPath(targetTime time.Time, dirPath string) (string, error) {
	year, month, day := targetTime.Date()
	// diaryDirPath := diaryDirPath()
	diaryPath := filepath.Join(
		dirPath,
		fmt.Sprintf("%02d", year),
		fmt.Sprintf("%02d", int(month)),
		fmt.Sprintf("%s.md", fmt.Sprintf("%02d", day)),
	)
	return diaryPath, nil
}

func targetTime(date string, before, after int) (time.Time, error) {
	now := time.Now()
	if date != "" {
		fmt.Println(date)
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

func makeFile(fPath string) error {
	err := ioutil.WriteFile(fPath, []byte(""), 0644)
	if err != nil {
		return fmt.Errorf("Failed make file. %v", err.Error())
	}
	return nil
}

func isFileExist(fPath string) bool {
	_, err := os.Stat(fPath)
	return err == nil || !os.IsNotExist(err)
}

func searchBeforeDate(startDate time.Time, dirPath string) (time.Time, error) {

	for i := 0; i < 30; i++ {
		date := startDate.AddDate(0, 0, -1*i)
		path, err := diaryPath(date, diaryDirPath())
		if err != nil {
			return time.Time{}, err
		}

		if isFileExist(path) {
			return date, nil
		}
	}
	return time.Time{}, fmt.Errorf("Not found before diary")
}

func openEditor(program string, args ...string) error {
	c := exec.Command(program, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
