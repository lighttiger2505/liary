package main

import (
	"fmt"
	"github.com/lighttiger2505/liary/internal"
	"os"
)

func Append(path string, val string, numLineBefore, numLineAfter int) error {
	// Make diary file
	err := internal.MakeFile(path)
	if err != nil {
		return err
	}

	// Open diary file
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("Failed append diary. %s", err.Error())
	}
	defer file.Close()

	// Append content
	appendBlackLine(file, numLineBefore)
	fmt.Fprintln(file, val)
	appendBlackLine(file, numLineAfter)

	return nil
}
func AppendCodeBlock(path string, val string, numLineBefore, numLineAfter int, lang string) error {
	// Make diary file
	err := internal.MakeFile(path)
	if err != nil {
		return err
	}

	// Open diary file
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("Failed append diary. %s", err.Error())
	}
	defer file.Close()

	// Append content
	appendBlackLine(file, numLineBefore)
	fmt.Fprintf(file, "```%s", lang)
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, val)
	fmt.Fprintf(file, "```")
	appendBlackLine(file, numLineAfter)

	return nil
}

func appendBlackLine(file *os.File, num int) {
	for i := 0; i < num; i++ {
		fmt.Fprintln(file, "")
	}
}
