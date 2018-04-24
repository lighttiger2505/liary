package main

import (
	"fmt"
	"os"
)

func Append(path string, val string) error {
	// Append content
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("Failed append diary. %s", err.Error())
	}
	defer file.Close()
	fmt.Fprintln(file, val)
	return nil
}
