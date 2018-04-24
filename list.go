package main

import "fmt"

func List() error {
	diaryDirPath := diaryDirPath()
	diaryPaths := dirWalk(diaryDirPath)
	for _, diaryPath := range diaryPaths {
		fmt.Println(diaryPath)
	}
	return nil
}
