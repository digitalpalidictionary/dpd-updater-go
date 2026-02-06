package system

import (
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
