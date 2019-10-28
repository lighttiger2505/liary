package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

var EditCommand = cli.Command{
	Name:      "edit",
	Aliases:   []string{"e"},
	Usage:     "edit diary",
	UsageText: "liary mv [command options...] <file suffix>",
	Action:    EditAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Open specified file",
		},
		cli.StringFlag{
			Name:  "date, d",
			Usage: "Open specified date diary",
		},
		cli.IntFlag{
			Name:  "before, b",
			Usage: "Open specified before diary by day",
		},
		cli.IntFlag{
			Name:  "after, a",
			Usage: "Open specified after diary by day",
		},
		cli.StringFlag{
			Name:  "subdir, s",
			Usage: "subdir",
		},
	},
}

func EditAction(c *cli.Context) error {
	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}

	var targetPath string
	file := c.String("file")
	if file != "" {
		p, err := getTargetPathWithFile(cfg.DiaryDir, file)
		if err != nil {
			return err
		}
		targetPath = p
	} else {
		diaryDir := cfg.DiaryDir
		subDir := c.String("subdir")
		if subDir != "" {
			diaryDir = filepath.Join(diaryDir, subDir)
		}

		p, err := getTargetPath(c, diaryDir)
		if err != nil {
			return err
		}
		targetPath = p
	}

	// Make directory
	targetDirPath := filepath.Dir(targetPath)
	if err := internal.MakeDir(targetDirPath); err != nil {
		return err
	}

	cmdArgs := []string{}
	if len(cfg.EditorOptions) > 0 {
		cmdArgs = append(cmdArgs, cfg.EditorOptions...)
	}
	cmdArgs = append(cmdArgs, targetPath)

	// Open text editor
	return internal.OpenEditor(cfg.Editor, cmdArgs...)
}

func getTargetPathWithFile(diaryDir, file string) (string, error) {
	var absSourcePath string
	if !filepath.IsAbs(file) {
		absSourcePath = filepath.Join(diaryDir, file)
	} else {
		absSourcePath = file
	}

	if !internal.IsFileExist(absSourcePath) {
		return "", fmt.Errorf("missing target file operand after, '%v'", absSourcePath)
	}
	return absSourcePath, nil
}

func getTargetPath(c *cli.Context, diaryDir string) (string, error) {
	// Getting time for target diary
	date := c.String("date")
	before := c.Int("before")
	after := c.Int("after")
	targetTime, err := internal.GetTargetTime(date, before, after)
	if err != nil {
		return "", err
	}

	// Getting diary path
	suffix := ""
	args := c.Args()
	if len(args) > 0 {
		suffix = suffixJoin(args[0])
	}
	targetPath, err := internal.DiaryPath(targetTime, diaryDir, suffix)
	if err != nil {
		return "", err
	}
	return targetPath, nil
}

func suffixJoin(val string) string {
	words := strings.Fields(val)
	return strings.Join(words, "_")
}
