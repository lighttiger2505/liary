package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandHomePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty path",
			input:    "",
			expected: "",
		},
		{
			name:     "Absolute path",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "Relative path",
			input:    "relative/path",
			expected: "relative/path",
		},
		{
			name:     "Home directory only",
			input:    "~",
			expected: homeDir,
		},
		{
			name:     "Home directory with subdirectory",
			input:    "~/Documents",
			expected: filepath.Join(homeDir, "Documents"),
		},
		{
			name:     "Home directory with nested path",
			input:    "~/Documents/diary",
			expected: filepath.Join(homeDir, "Documents", "diary"),
		},
		{
			name:     "HOME environment variable",
			input:    "$HOME/Documents",
			expected: filepath.Join(homeDir, "Documents"),
		},
		{
			name:     "HOME environment variable with braces",
			input:    "${HOME}/Documents",
			expected: filepath.Join(homeDir, "Documents"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandHomePath(tt.input)
			if result != tt.expected {
				t.Errorf("expandHomePath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConfigLoadExpandsPaths(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yml")

	// Test config content with home directory paths
	configContent := `diarydir: ~/diary
editor: vim
editoroptions: []
workspaces:
  default: ~/diary
  work: $HOME/work-diary
grepcmd: grep -nH ${PATTERN} ${FILES}
`

	err := os.WriteFile(configPath, []byte(configContent), 0666)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Temporarily replace the global config file path
	originalConfigFilePath := configFilePath
	configFilePath = configPath
	defer func() {
		configFilePath = originalConfigFilePath
	}()

	// Load config
	cfg := newConfig()
	err = cfg.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	// Check that paths were expanded
	expectedDiaryDir := filepath.Join(homeDir, "diary")
	if cfg.DiaryDir != expectedDiaryDir {
		t.Errorf("DiaryDir not expanded: got %q, want %q", cfg.DiaryDir, expectedDiaryDir)
	}

	expectedDefaultWorkspace := filepath.Join(homeDir, "diary")
	if cfg.WorkSpaces["default"] != expectedDefaultWorkspace {
		t.Errorf("Default workspace not expanded: got %q, want %q", cfg.WorkSpaces["default"], expectedDefaultWorkspace)
	}

	expectedWorkWorkspace := filepath.Join(homeDir, "work-diary")
	if cfg.WorkSpaces["work"] != expectedWorkWorkspace {
		t.Errorf("Work workspace not expanded: got %q, want %q", cfg.WorkSpaces["work"], expectedWorkWorkspace)
	}
}
