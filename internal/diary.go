package internal

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

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

func DiaryPath(targetTime time.Time, dirPath string, suffix string) (string, error) {
	_, _, day := targetTime.Date()
	var filename string
	if suffix != "" {
		filename = fmt.Sprintf("%s-%s.md", fmt.Sprintf("%02d", day), suffix)
	} else {
		filename = fmt.Sprintf("%s.md", fmt.Sprintf("%02d", day))
	}

	diaryPath := filepath.Join(
		MonthPath(targetTime, dirPath),
		filename,
	)
	return diaryPath, nil
}

func GetDiaryList(diaryDirPath string, isAll bool, isFullPath bool, dateRangeString string) ([]string, error) {
	diaryList := FilterMarkdown(Walk(diaryDirPath))
	if !isAll {
		filterdList, err := filterDateRange(diaryList, diaryDirPath, dateRangeString)
		if err != nil {
			return nil, err
		}
		diaryList = filterdList
	}

	showPaths := diaryList
	if !isFullPath {
		showPaths = []string{}
		for _, diaryPath := range diaryList {
			showPaths = append(showPaths, strings.TrimPrefix(diaryPath, diaryDirPath+"/"))
		}
	}
	return showPaths, nil
}

func filterDateRange(base []string, diaryDirPath string, dateRangeString string) ([]string, error) {
	year, month, day, err := ParseDate(dateRangeString)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	start := now.AddDate(-1*year, -1*month, -1*day)
	dateRange := GetDateRange(start, now)

	targetPaths := []string{}
	for _, d := range dateRange {
		targetPaths = append(targetPaths, DayPath(d, diaryDirPath))
	}

	// Show filtering list
	filteredPaths := []string{}
	for _, p := range base {
		for _, targetPath := range targetPaths {
			if strings.HasPrefix(p, targetPath) {
				filteredPaths = append(filteredPaths, p)
			}
		}
	}

	return filteredPaths, nil
}
