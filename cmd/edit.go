package cmd

import (
	"fmt"
	"os"
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
	},
}

func EditAction(c *cli.Context) error {
	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}

	fmt.Println(c.GlobalFlagNames())
	workspace := cfg.DiaryDir
	workspaceFlag := c.GlobalString("workspace")
	if workspaceFlag != "" {
		w, err := getWorkSpace(cfg.WorkSpaces, workspaceFlag)
		if err != nil {
			return err
		}
		workspace = w
	}

	var targetPath string
	file := c.String("file")
	if file != "" {
		p, err := getTargetPathWithFile(workspace, file)
		if err != nil {
			return err
		}
		targetPath = p
	} else {
		p, err := getTargetPath(c, workspace)
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

func getWorkSpace(workspaces map[string]string, name string) (string, error) {
	if len(workspaces) == 0 {
		return "", fmt.Errorf("Not set workspace")
	}

	workspace, ok := workspaces[name]
	if !ok {
		return "", fmt.Errorf("Not found workspace, %s", name)
	}
	if !internal.IsFileExist(workspace) {
		return "", fmt.Errorf("No such directory, %s", workspace)
	}

	f, err := os.Open(workspace)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fs, err := os.Stat(workspace)
	if err != nil {
		return "", err
	}
	if !fs.IsDir() {
		return "", fmt.Errorf("Workspace not directory, %s", workspace)
	}

	return workspace, nil
}
