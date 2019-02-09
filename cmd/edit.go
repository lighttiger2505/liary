package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

func EditAction(c *cli.Context) error {
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
	targetPath, err := internal.DiaryPath(targetTime, internal.DiaryDirPath(), suffix)
	if err != nil {
		return err
	}

	// Show diary file path
	path := c.Bool("path")
	if path {
		fmt.Println(targetPath)
		return nil
	}

	// Make directory
	targetDirPath := filepath.Dir(targetPath)
	if err := internal.MakeDir(targetDirPath); err != nil {
		return err
	}

	// Open text editor
	return edit(targetPath)
}

func edit(path string) error {
	// Open text editor
	editorEnv := os.Getenv("EDITOR")
	if editorEnv == "" {
		editorEnv = "vim"
	}
	if err := openEditor(editorEnv, path); err != nil {
		return fmt.Errorf("Failed open editor. %s", err.Error())
	}
	return nil
}

func suffixJoin(val string) string {
	words := strings.Fields(val)
	return strings.Join(words, "_")
}

func openEditor(program string, args ...string) error {
	c := exec.Command(program, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
