package installer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeleteFolders(t *testing.T) {
	tmpDir := t.TempDir()

	// Create structure:
	// tmpDir/folder1/dpd.ifo
	// tmpDir/folder2/dpd.ifo
	// tmpDir/folder3/ (should be kept)

	dir1 := filepath.Join(tmpDir, "folder1")
	dir2 := filepath.Join(tmpDir, "folder2")
	dir3 := filepath.Join(tmpDir, "folder3")

	os.Mkdir(dir1, 0755)
	os.Mkdir(dir2, 0755)
	os.Mkdir(dir3, 0755)

	// Create dummy ifo files
	os.WriteFile(filepath.Join(dir1, "dpd.ifo"), []byte("data"), 0644)
	os.WriteFile(filepath.Join(dir2, "dpd.ifo"), []byte("data"), 0644)

	// Inputs to DeleteFolders are paths to dpd.ifo files
	pathsToDelete := []string{
		filepath.Join(dir1, "dpd.ifo"),
		filepath.Join(dir2, "dpd.ifo"),
	}

	err := DeleteFolders(pathsToDelete)
	if err != nil {
		t.Fatalf("DeleteFolders failed: %v", err)
	}

	// Verify folder1 and folder2 are gone
	if _, err := os.Stat(dir1); !os.IsNotExist(err) {
		t.Error("folder1 should have been deleted")
	}
	if _, err := os.Stat(dir2); !os.IsNotExist(err) {
		t.Error("folder2 should have been deleted")
	}

	// Verify folder3 still exists
	if _, err := os.Stat(dir3); os.IsNotExist(err) {
		t.Error("folder3 should NOT have been deleted")
	}
}

func TestDeleteFoldersSafeChecks(t *testing.T) {
	// Test unsafe paths
	unsafePaths := []string{
		"/dpd.ifo",
	}

	err := DeleteFolders(unsafePaths)
	if err == nil {
		t.Error("Expected error for root deletion attempt")
	}
}
