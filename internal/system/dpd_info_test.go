package system

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseIFO(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create a valid dpd.ifo file
	validContent := `version=3.0.0
bookname=Digital Pāḷi Dictionary (pi-en)
author=Bodhirasa
date=2026-02-06T14:40:12
description=Test Description`
	validPath := filepath.Join(tmpDir, "dpd.ifo")
	if err := os.WriteFile(validPath, []byte(validContent), 0644); err != nil {
		t.Fatalf("Failed to write valid IFO file: %v", err)
	}

	// Create a malformed dpd.ifo file (missing date)
	malformedContent := `version=3.0.0
bookname=Digital Pāḷi Dictionary (pi-en)
description=Test Description`
	malformedPath := filepath.Join(tmpDir, "malformed.ifo")
	if err := os.WriteFile(malformedPath, []byte(malformedContent), 0644); err != nil {
		t.Fatalf("Failed to write malformed IFO file: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		wantErr  bool
		wantDate string
	}{
		{
			name:     "Valid IFO",
			path:     validPath,
			wantErr:  false,
			wantDate: "2026-02-06T14:40:12",
		},
		{
			name:    "Malformed IFO",
			path:    malformedPath,
			wantErr: true,
		},
		{
			name:    "Non-existent File",
			path:    filepath.Join(tmpDir, "nonexistent.ifo"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ParseIFO(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIFO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if info.Path != tt.path {
					t.Errorf("ParseIFO() Path = %v, want %v", info.Path, tt.path)
				}
				// Parse expected date for comparison
				expectedTime, _ := time.Parse("2006-01-02T15:04:05", tt.wantDate)
				if !info.Date.Equal(expectedTime) {
					t.Errorf("ParseIFO() Date = %v, want %v", info.Date, expectedTime)
				}
			}
		})
	}
}
