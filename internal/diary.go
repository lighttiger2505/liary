package internal

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

func DiaryDirPath() string {
	configPath := getXDGConfigPath(runtime.GOOS)
	diaryDirPath := filepath.Join(configPath, "_post")
	return diaryDirPath
}

func YearPath(targetTime time.Time, dirPath string) string {
	year, _, _ := targetTime.Date()
	result := filepath.Join(
		dirPath,
		fmt.Sprintf("%02d", year),
	)
	return result
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

func DiaryPath(targetTime time.Time, suffix string) (string, error) {
	_, _, day := targetTime.Date()
	var filename string
	if suffix != "" {
		filename = fmt.Sprintf("%s-%s.md", fmt.Sprintf("%02d", day), suffix)
	} else {
		filename = fmt.Sprintf("%s.md", fmt.Sprintf("%02d", day))
	}

	diaryPath := filepath.Join(
		MonthPath(targetTime, DiaryDirPath()),
		filename,
	)
	return diaryPath, nil
}
