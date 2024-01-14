package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func GetAppendValue(args []string) (string, error) {
	var val string
	if terminal.IsTerminal(0) {
		if len(args) > 0 {
			val = args[0]
		} else {
			val = ""
		}
	} else {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("Failed make diary file. %s", err.Error())
		}
		val = string(b)
	}
	return val, nil
}

func OpenEditor(program string, args ...string) error {
	cmdargs := strings.Join(args, " ")
	command := program + " " + cmdargs

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GrepFiles(command, pattern string, files ...string) error {
	var args []string
	for _, file := range files {
		args = append(args, shellquote(file))
	}
	cmdargs := strings.Join(args, " ")

	hasEnv := false
	command = os.Expand(command, func(s string) string {
		hasEnv = true
		switch s {
		case "FILES":
			return cmdargs
		case "PATTERN":
			return pattern
		}
		return os.Getenv(s)
	})
	if !hasEnv {
		command += " " + cmdargs
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func shellquote(s string) string {
	return fmt.Sprintf("%q", s)
}
