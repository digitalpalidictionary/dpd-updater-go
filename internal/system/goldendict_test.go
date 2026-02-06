package system

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetGoldenDictConfigPath(t *testing.T) {
	t.Run("Returns correct path based on OS", func(t *testing.T) {
		tempHome := t.TempDir()

		// Mock environment variables to point to the temp directory
		// os.UserConfigDir() relies on these
		if runtime.GOOS == "windows" {
			t.Setenv("APPDATA", tempHome)
		} else {
			t.Setenv("HOME", tempHome)
			t.Setenv("XDG_CONFIG_HOME", filepath.Join(tempHome, ".config"))
		}

		got, err := GetGoldenDictConfigPath()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		var expected string
		switch runtime.GOOS {
		case "windows":
			// Windows: %APPDATA%\GoldenDict\config
			expected = filepath.Join(tempHome, "GoldenDict", "config")
		case "darwin":
			// macOS: ~/Library/Application Support/GoldenDict/config
			// NOTE: os.UserConfigDir() on macOS returns ~/Library/Application Support
			expected = filepath.Join(tempHome, "Library", "Application Support", "GoldenDict", "config")
		default:
			// Linux/Unix: ~/.config/goldendict/config
			expected = filepath.Join(tempHome, ".config", "goldendict", "config")
		}

		if got != expected {
			t.Errorf("Expected %s, got %s", expected, got)
		}
	})
}

func TestParseGoldenDictPaths(t *testing.T) {
	t.Run("Parses valid config correctly", func(t *testing.T) {
		content := `
<config>
  <paths>
    <path recursive="true" enabled="true">/path/to/recursive</path>
    <path recursive="0" enabled="true">/path/to/flat</path>
    <path recursive="true" enabled="false">/path/to/disabled</path>
    <path recursive="1">/path/to/no-enabled-attr</path>
  </paths>
</config>`
		tmpFile := filepath.Join(t.TempDir(), "config")
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		paths, err := ParseGoldenDictPaths(tmpFile)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(paths) != 3 {
			t.Errorf("Expected 3 paths, got %d", len(paths))
		}

		// Check first path
		if len(paths) > 0 {
			if paths[0].Path != "/path/to/recursive" || !paths[0].Recursive {
				t.Errorf("Mismatch in first path: %+v", paths[0])
			}
		}

		// Check third path (the one with missing enabled attr)
		if len(paths) > 2 {
			if paths[2].Path != "/path/to/no-enabled-attr" || !paths[2].Recursive || !paths[2].Enabled {
				t.Errorf("Mismatch in third path: %+v", paths[2])
			}
		}
	})

	t.Run("Handles missing file gracefully", func(t *testing.T) {
		_, err := ParseGoldenDictPaths("/non/existent/path")
		if err == nil {
			t.Error("Expected error for missing file, got nil")
		}
	})
}

func TestAnalyzeGoldenDictPaths(t *testing.T) {
	tests := []struct {
		name     string
		paths    []GDPath
		expected string
	}{
		{
			name:     "Single recursive path",
			paths:    []GDPath{{Path: "/home/user/dicts", Recursive: true, Enabled: true}},
			expected: "/home/user/dicts",
		},
		{
			name:     "Single non-recursive path",
			paths:    []GDPath{{Path: "/home/user/dicts", Recursive: false, Enabled: true}},
			expected: "",
		},
		{
			name: "Multiple paths with common parent",
			paths: []GDPath{
				{Path: "/home/user/dicts/Pali", Recursive: true, Enabled: true},
				{Path: "/home/user/dicts/Sanskrit", Recursive: true, Enabled: true},
			},
			expected: "/home/user/dicts",
		},
		{
			name: "Multiple paths with no common parent",
			paths: []GDPath{
				{Path: "/home/user/dicts", Recursive: true, Enabled: true},
				{Path: "/opt/dicts", Recursive: true, Enabled: true},
			},
			expected: "",
		},
		{
			name:     "Empty paths",
			paths:    []GDPath{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AnalyzeGoldenDictPaths(tt.paths)
			if got != tt.expected {
				t.Errorf("AnalyzeGoldenDictPaths() = %v, want %v", got, tt.expected)
			}
		})
	}
}
