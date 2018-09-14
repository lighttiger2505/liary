package internal

import (
	"fmt"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

func DiaryDirPath() string {
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
