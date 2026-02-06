package system

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func ScanForVersion(gdPath string) (string, error) {
	if gdPath == "" {
		return "unknown", nil
	}

	// Look for dpd folder
	dpdFolder := filepath.Join(gdPath, "dpd")
	if _, err := os.Stat(dpdFolder); os.IsNotExist(err) {
		return "unknown", nil
	}

	// Scan .ifo files in the dpd folder
	files, err := os.ReadDir(dpdFolder)
	if err != nil {
		return "unknown", err
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".ifo") {
			ifoPath := filepath.Join(dpdFolder, f.Name())
			file, err := os.Open(ifoPath)
			if err != nil {
				continue
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "date=") {
					return strings.TrimPrefix(line, "date="), nil
				}
			}
		}
	}

	return "installed", nil
}

func ValidateGoldenDictPath(path string) (bool, string) {
	info, err := os.Stat(path)
	if err != nil {
		return false, "Path does not exist"
	}
	if !info.IsDir() {
		return false, "Path is not a directory"
	}

	// Check if user selected a DPD subfolder instead of parent
	dpdFolders := []string{"dpd", "dpd-grammar", "dpd-deconstructor", "dpd-variants"}
	name := filepath.Base(path)
	for _, f := range dpdFolders {
		if strings.EqualFold(name, f) {
			return true, "Warning: You selected a subfolder. Using parent is recommended."
		}
	}

	// Check for dictionary indicators
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, "Error reading directory"
	}

	if len(entries) == 0 {
		return false, "Directory is empty"
	}

	hasDicts := false
	for _, entry := range entries {
		if entry.IsDir() {
			hasDicts = true
			break
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".ifo" || ext == ".dsl" || ext == ".zip" || ext == ".dz" {
			hasDicts = true
			break
		}
	}

	if !hasDicts {
		return false, "No dictionary files or folders found"
	}

	return true, "Valid GoldenDict folder"
}