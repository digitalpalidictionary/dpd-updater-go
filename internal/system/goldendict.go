package system

import (
	"encoding/xml"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type GoldenDictManager struct {
	ProcessNames []string
}

func NewGoldenDictManager() *GoldenDictManager {
	return &GoldenDictManager{
		ProcessNames: []string{"goldendict", "GoldenDict", "goldendict.exe", "GoldenDict.exe"},
	}
}

func (gm *GoldenDictManager) IsRunning() (bool, error) {
	processes, err := process.Processes()
	if err != nil {
		return false, err
	}

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		for _, pn := range gm.ProcessNames {
			if strings.EqualFold(name, pn) {
				return true, nil
			}
		}
	}

	return false, nil
}

func (gm *GoldenDictManager) Close(timeout time.Duration) error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		match := false
		for _, pn := range gm.ProcessNames {
			if strings.EqualFold(name, pn) {
				match = true
				break
			}
		}

		if match {
			p.Terminate()
		}
	}

	start := time.Now()
	for time.Since(start) < timeout {
		running, _ := gm.IsRunning()
		if !running {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Force kill if still running
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		match := false
		for _, pn := range gm.ProcessNames {
			if strings.EqualFold(name, pn) {
				match = true
				break
			}
		}

		if match {
			p.Kill()
		}
	}

	return nil
}

func (gm *GoldenDictManager) Reopen() error {
	exe, err := gm.FindExecutable()
	if err != nil {
		return err
	}

	cmd := exec.Command(exe)
	// Platform specific detachment
	return cmd.Start()
}

func (gm *GoldenDictManager) FindExecutable() (string, error) {
	// 1. Check PATH
	if path, err := exec.LookPath("goldendict"); err == nil {
		return path, nil
	}

	// 2. Check common locations
	var commonPaths []string
	switch runtime.GOOS {
	case "windows":
		commonPaths = []string{
			"C:\\Program Files\\GoldenDict\\GoldenDict.exe",
			"C:\\Program Files (x86)\\GoldenDict\\GoldenDict.exe",
			filepath.Join(os.Getenv("LOCALAPPDATA"), "GoldenDict", "GoldenDict.exe"),
		}
	case "darwin":
		commonPaths = []string{
			"/Applications/GoldenDict.app/Contents/MacOS/GoldenDict",
			filepath.Join(os.Getenv("HOME"), "Applications/GoldenDict.app/Contents/MacOS/GoldenDict"),
		}
	case "linux":
		commonPaths = []string{
			"/usr/bin/goldendict",
			"/usr/local/bin/goldendict",
			"/opt/goldendict/goldendict",
		}
	}

	for _, p := range commonPaths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", os.ErrNotExist
}

// GetGoldenDictConfigPath returns the platform-specific path to the GoldenDict config file.
func GetGoldenDictConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	folderName := "GoldenDict"
	if runtime.GOOS == "linux" {
		folderName = "goldendict"
	}

	return filepath.Join(configDir, folderName, "config"), nil
}

// GDPath represents a dictionary path entry in GoldenDict config
type GDPath struct {
	Path      string
	Recursive bool
	Enabled   bool
}

// xml structures for parsing
type gdConfigXML struct {
	Paths struct {
		Path []struct {
			Recursive string `xml:"recursive,attr"`
			Enabled   string `xml:"enabled,attr"`
			Value     string `xml:",chardata"`
		} `xml:"path"`
	} `xml:"paths"`
}

// ParseGoldenDictPaths reads the config file and returns a list of enabled dictionary paths.
func ParseGoldenDictPaths(configPath string) ([]GDPath, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config gdConfigXML
	if err := xml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	var results []GDPath
	for _, p := range config.Paths.Path {
		enabled := isTrue(p.Enabled)
		if !enabled {
			continue
		}

		results = append(results, GDPath{
			Path:      p.Value,
			Recursive: isTrue(p.Recursive),
			Enabled:   enabled,
		})
	}

	return results, nil
}

func isTrue(s string) bool {
	s = strings.ToLower(s)
	return s == "true" || s == "1" || s == "yes" || s == "on"
}

// AnalyzeGoldenDictPaths analyzes the list of paths and returns a suggested installation folder.
// Returns an empty string if no suitable suggestion is found.
func AnalyzeGoldenDictPaths(paths []GDPath) string {
	if len(paths) == 0 {
		return ""
	}

	if len(paths) == 1 {
		if paths[0].Recursive {
			return filepath.Clean(paths[0].Path)
		}
		return ""
	}

	common := filepath.Clean(paths[0].Path)
	for _, p := range paths[1:] {
		path := filepath.Clean(p.Path)
		// Iteratively move up common until it contains path
		for {
			rel, err := filepath.Rel(common, path)
			if err == nil && !strings.HasPrefix(rel, "..") {
				break // common is parent of path
			}

			parent := filepath.Dir(common)
			if parent == common {
				// We reached the root.
				// If we are at root, we stop here.
				common = parent
				break
			}
			common = parent
		}
	}

	// If the common path is the root directory, we consider it "no common parent"
	// because suggesting the root of the drive is rarely desired.
	if filepath.Dir(common) == common || common == "." || common == "" {
		return ""
	}

	return common
}
