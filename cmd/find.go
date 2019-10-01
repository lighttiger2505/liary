package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

var FindCommand = cli.Command{
	Name:    "find",
	Aliases: []string{"f"},
	Usage:   "find and edit diary",
	Action:  FindAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "range, r",
			Usage: "relative date range",
			Value: DefaultDateRange,
		},
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "show all diary's",
		},
	},
}

func FindAction(c *cli.Context) error {
	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}

	p, err := fuzzyFindDiary(c, cfg)
	if err != nil {
		return err
	}

	cmdArgs := []string{}
	if len(cfg.EditorOptions) > 0 {
		cmdArgs = append(cmdArgs, cfg.EditorOptions...)
	}
	cmdArgs = append(cmdArgs, p)

	// Open text editor
	return internal.OpenEditor(cfg.Editor, cmdArgs...)
}

func fuzzyFindDiary(c *cli.Context, cfg *internal.Config) (string, error) {
	paths, err := GetDiaryList(cfg.DiaryDir, false, false, c.String("range"))
	if err != nil {
		return "", err
	}

	index, err := fuzzyfinder.Find(
		paths,
		func(i int) string {
			return paths[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}

			f, err := os.Open(filepath.Join(cfg.DiaryDir, paths[i]))
			if err != nil {
				return "file open error..."
			}
			defer f.Close()

			b, err := ioutil.ReadAll(f)
			if err != nil {
				return "read error..."
			}
			return string(b)
		}),
	)

	if err != nil {
		if err.Error() == fuzzyfinder.ErrAbort.Error() {
			return "", errors.New("interrupted")
		}
		return "", err
	}
	return filepath.Join(cfg.DiaryDir, paths[index]), nil
}
