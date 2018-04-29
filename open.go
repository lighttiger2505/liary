package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func Open(path string) error {
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

func makeFile(fPath string) error {
	err := ioutil.WriteFile(fPath, []byte(""), 0644)
	if err != nil {
		return fmt.Errorf("Failed make file. %v", err.Error())
	}
	return nil
}

func openEditor(program string, args ...string) error {
	c := exec.Command(program, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
