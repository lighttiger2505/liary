package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

func ListAction(c *cli.Context) error {
	// Show diary file list
	return list()
}

func list() error {
	diaryDirPath := internal.DiaryDirPath()
	diaryPaths := dirWalk(diaryDirPath)
	for _, diaryPath := range diaryPaths {
		fmt.Println(diaryPath)
	}
	return nil
}

func ListTargetDate(date time.Time) error {
	monthPath := internal.MonthPath(date, internal.DiaryDirPath())
	dayPath := internal.DayPath(date, internal.DiaryDirPath())
	diaryPaths := dirWalk(monthPath)
	for _, diaryPath := range diaryPaths {
		if strings.HasPrefix(diaryPath, dayPath) {
			fmt.Println(diaryPath)
		}
	}
	return nil
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
