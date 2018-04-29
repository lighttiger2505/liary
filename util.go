package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh/terminal"
)

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

func DiaryPath(targetTime time.Time, dirPath string, suffix string) (string, error) {
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

func TargetTime(date string, before, after int) (time.Time, error) {
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
		path, err := DiaryPath(date, diaryDirPath(), "")
		if err != nil {
			return time.Time{}, err
		}

		if isFileExist(path) {
			return date, nil
		}
	}
	return time.Time{}, fmt.Errorf("Not found before diary")
}

func GetAppendValue(args []string) (string, error) {
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
