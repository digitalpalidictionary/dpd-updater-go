package system

import (
	"os"
	"path/filepath"
	"strings"
)

func ScanForVersion(gdPath string) (string, error) {
	if gdPath == "" {
		return "unknown", nil
	}

	// Look for dpd.ifo or similar metadata files
	// In the original, it looks for specific folders and files
	// For now, let's implement a simple version check if a version file exists
	// or scan the .ifo files for version strings.

	versionFile := filepath.Join(gdPath, "dpd", "version.txt")
	if _, err := os.Stat(versionFile); err == nil {
		data, err := os.ReadFile(versionFile)
		if err == nil {
			return strings.TrimSpace(string(data)), nil
		}
	}

	return "unknown", nil
}

func ValidateGoldenDictPath(path string) (bool, string) {
	info, err := os.Stat(path)
	if err != nil {
		return false, "Path does not exist"
	}
	if !info.IsDir() {
		return false, "Path is not a directory"
	}

	// Check for DPD subfolders
	dpdFolders := []string{"dpd", "dpd-grammar", "dpd-deconstructor", "dpd-variants"}
	name := filepath.Base(path)
	for _, f := range dpdFolders {
		if strings.EqualFold(name, f) {
			return true, "Warning: You selected a subfolder. Using parent is recommended."
		}
	}

	return true, "Valid GoldenDict folder"
}
