package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func IsFileExist(fPath string) bool {
	_, err := os.Stat(fPath)
	return err == nil || !os.IsNotExist(err)
}

func MakeFile(fPath string) error {
	if !IsFileExist(fPath) {
		err := os.WriteFile(fPath, []byte(""), 0644)
		if err != nil {
			return fmt.Errorf("Failed make diary file. %v", err.Error())
		}
	}
	return nil
}

func MakeDir(dPath string) error {
	if !IsFileExist(dPath) {
		if err := os.MkdirAll(dPath, 0755); err != nil {
			return fmt.Errorf("Failed make diary dir. %s", err.Error())
		}
	}
	return nil
}

const APP_NAME = "liary"

func getXDGConfigPath(goos string) string {
	var dir string
	if goos == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", APP_NAME)
		}
		dir = filepath.Join(dir, "lab")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", APP_NAME)
	}
	return dir
}

func Walk(dir string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, Walk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}

func FilterMarkdown(files []string) []string {
	var newfiles []string
	for _, file := range files {
		if strings.HasSuffix(file, ".md") {
			newfiles = append(newfiles, file)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(newfiles)))
	return newfiles
}
