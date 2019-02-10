package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

func AppendAction(c *cli.Context) error {
	// Getting time for target diary
	date := c.String("date")
	before := c.Int("before")
	after := c.Int("after")
	targetTime, err := internal.GetTargetTime(date, before, after)
	if err != nil {
		return err
	}

	// Getting diary path
	suffix := suffixJoin(c.String("suffix"))
	targetPath, err := internal.DiaryPath(targetTime, suffix)
	if err != nil {
		return err
	}

	// Make directory
	targetDirPath := filepath.Dir(targetPath)
	if err := internal.MakeDir(targetDirPath); err != nil {
		return err
	}

	appendVal, err := internal.GetAppendValue(c.Args())
	if err != nil {
		return fmt.Errorf("Failed get append value %s", err)
	}

	code := c.Bool("code")
	lang := c.String("language")
	numLineBefore := c.Int("before-append")
	numLineAfter := c.Int("after-append")
	if appendVal != "" {
		if code {
			return appendCodeBlock(targetPath, appendVal, numLineBefore, numLineAfter, lang)
		}
		return appendText(targetPath, appendVal, numLineBefore, numLineAfter)
	}
	return nil
}

func appendText(path string, val string, numLineBefore, numLineAfter int) error {
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
func appendCodeBlock(path string, val string, numLineBefore, numLineAfter int, lang string) error {
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
