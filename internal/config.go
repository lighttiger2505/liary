package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	yaml "gopkg.in/yaml.v2"
)

var configFilePath = filepath.Join(getXDGConfigPath(runtime.GOOS), "config.yml")

type Config struct {
	DiaryDir      string   `yaml:"diarydir"`
	Editor        string   `yaml:"editor"`
	EditorOptions []string `yaml:"editoroptions"`
	GrepCmd       string   `yaml:"grepcmd"`
}

func GetConfig() (*Config, error) {
	cfg := newConfig()
	if err := cfg.Load(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Path() string {
	return configFilePath
}

func (c *Config) Read() (string, error) {
	if err := os.MkdirAll(filepath.Dir(configFilePath), 0700); err != nil {
		return "", fmt.Errorf("cannot create directory, %s", err)
	}

	if !IsFileExist(configFilePath) {
		_, err := os.Create(configFilePath)
		if err != nil {
			return "", fmt.Errorf("cannot create config, %s", err.Error())
		}
	}

	file, err := os.OpenFile(configFilePath, os.O_RDONLY, 0666)
	if err != nil {
		return "", fmt.Errorf("cannot open config, %s", err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("cannot read config, %s", err)
	}

	return string(b), nil
}

func (c *Config) Load() error {
	if err := os.MkdirAll(filepath.Dir(configFilePath), 0700); err != nil {
		return fmt.Errorf("cannot create directory, %s", err)
	}

	if !IsFileExist(configFilePath) {
		if err := createNewConfig(); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(configFilePath, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("cannot open config, %s", err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("cannot read config, %s", err)
	}

	if err = yaml.Unmarshal(b, c); err != nil {
		return fmt.Errorf("failed unmarshal yaml. \nError: %s \nBuffer: %s", err, string(b))
	}
	return nil
}

func (c *Config) Save() error {
	file, err := os.OpenFile(configFilePath, os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("cannot open file, %s", err)
	}
	defer file.Close()

	out, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("Failed marshal config. Error: %v", err)
	}

	if _, err = io.WriteString(file, string(out)); err != nil {
		return fmt.Errorf("Failed write config file. Error: %s", err)
	}
	return nil
}

func newConfig() *Config {
	cfg := &Config{}
	return cfg
}

func createNewConfig() error {
	// Create new config file
	_, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("cannot create config, %s", err.Error())
	}

	// Add default settings
	cfg := newConfig()

	configPath := getXDGConfigPath(runtime.GOOS)
	diaryDirPath := filepath.Join(configPath, "_post")
	cfg.DiaryDir = diaryDirPath

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	cfg.Editor = editor

	cfg.GrepCmd = "grep -nH ${PATTERN} ${FILES}"

	cfg.Save()
	return nil
}
