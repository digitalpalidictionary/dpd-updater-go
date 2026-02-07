package system

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindAllDPDInstances(t *testing.T) {
	tmpDir := t.TempDir()

	// Structure:
	// tmpDir/
	//   valid1/dpd.ifo (2026-02-06)
	//   valid2/dpd.ifo (2026-01-01)
	//   empty/
	//   malformed/dpd.ifo (no date)
	//   other/something.txt

	// 1. Valid Instance 1 (Newer)
	valid1Dir := filepath.Join(tmpDir, "valid1")
	os.Mkdir(valid1Dir, 0755)
	valid1Content := `version=3.0.0
bookname=Digital Pāḷi Dictionary (pi-en)
author=Bodhirasa
date=2026-02-06T14:40:12`
	os.WriteFile(filepath.Join(valid1Dir, "dpd.ifo"), []byte(valid1Content), 0644)

	// 2. Valid Instance 2 (Older)
	valid2Dir := filepath.Join(tmpDir, "valid2")
	os.Mkdir(valid2Dir, 0755)
	valid2Content := `version=2.0.0
bookname=Digital Pāḷi Dictionary (pi-en)
author=Bodhirasa
date=2026-01-01T10:00:00`
	os.WriteFile(filepath.Join(valid2Dir, "dpd.ifo"), []byte(valid2Content), 0644)

	// 3. Malformed Instance
	malformedDir := filepath.Join(tmpDir, "malformed")
	os.Mkdir(malformedDir, 0755)
	malformedContent := `version=1.0.0
bookname=DPD`
	os.WriteFile(filepath.Join(malformedDir, "dpd.ifo"), []byte(malformedContent), 0644)

	// 4. Other file
	otherDir := filepath.Join(tmpDir, "other")
	os.Mkdir(otherDir, 0755)
	os.WriteFile(filepath.Join(otherDir, "something.txt"), []byte("hello"), 0644)

	// Run Test
	instances, err := FindAllDPDInstances(tmpDir)
	if err != nil {
		t.Fatalf("FindAllDPDInstances failed: %v", err)
	}

	if len(instances) != 2 {
		t.Errorf("Expected 2 instances, got %d", len(instances))
	}

	// Verify we found the correct paths
	foundPaths := make(map[string]bool)
	for _, inst := range instances {
		foundPaths[inst.Path] = true
	}

	if !foundPaths[filepath.Join(valid1Dir, "dpd.ifo")] {
		t.Error("Did not find valid1/dpd.ifo")
	}
	if !foundPaths[filepath.Join(valid2Dir, "dpd.ifo")] {
		t.Error("Did not find valid2/dpd.ifo")
	}
}
