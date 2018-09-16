package internal

import (
	"fmt"
	"io/ioutil"
	"os"
)

func isFileExist(fPath string) bool {
	_, err := os.Stat(fPath)
	return err == nil || !os.IsNotExist(err)
}

func MakeFile(fPath string) error {
	if !isFileExist(fPath) {
		err := ioutil.WriteFile(fPath, []byte(""), 0644)
		if err != nil {
			return fmt.Errorf("Failed make diary file. %v", err.Error())
		}
	}
	return nil
}

func MakeDir(dPath string) error {
	if !isFileExist(dPath) {
		if err := os.MkdirAll(dPath, 0755); err != nil {
			return fmt.Errorf("Failed make diary dir. %s", err.Error())
		}
	}
	return nil
}
