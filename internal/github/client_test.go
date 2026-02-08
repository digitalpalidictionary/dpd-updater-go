package github

import (
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		current string
		latest  string
		want    int // -1 if current < latest, 1 if current > latest, 0 if equal
	}{
		// Date-based versions
		{"2026-02-07", "2026-02-08", -1},
		{"2026-02-08", "2026-02-07", 1},
		{"2026-02-07", "2026-02-07", 0},

		// ISO 8601 with time vs Date-only
		{"2026-02-07T10:21:14", "2026-02-08", -1},
		{"2026-02-07T10:21:14", "2026-02-07", 0},

		// Real-world case: IFO date vs GitHub semver with date as patch
		{"2025-07-09", "v0.3.20260202", -1}, // Should detect update available
		{"2026-02-02", "v0.3.20260202", 0},  // Same date
		{"2026-02-03", "v0.3.20260202", 1},  // Installed is newer (unreleased version)

		// Semantic-like or string fallback
		{"v1.0.0", "v1.0.1", -1},
		{"v1.1.0", "v1.0.9", 1},
		{"unknown", "2026-02-07", -1},
	}

	for _, tt := range tests {
		t.Run(tt.current+"_vs_"+tt.latest, func(t *testing.T) {
			got := CompareVersions(tt.current, tt.latest)
			if got != tt.want {
				t.Errorf("CompareVersions(%s, %s) = %v, want %v", tt.current, tt.latest, got, tt.want)
			}
		})
	}
}
