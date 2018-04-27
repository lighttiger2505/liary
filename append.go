package main

import (
	"fmt"
	"os"
)

func Append(path string, val string) error {
	// Make diary file
	if !isFileExist(path) {
		if err := makeFile(path); err != nil {
			return fmt.Errorf("Failed make diary file. %s", err.Error())
		}
	}

	// Append content
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("Failed append diary. %s", err.Error())
	}
	defer file.Close()
	fmt.Fprintln(file, val)
	return nil
}
