package main

import (
	"fmt"
	"github.com/lighttiger2505/liary/internal"
	"os"
)

func Append(path string, val string) error {
	// Make diary file
	if !internal.IsFileExist(path) {
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
