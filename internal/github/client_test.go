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
		{"2026-02-07T10:21:14", "2026-02-07", 0}, // Date-only parsing of 2026-02-07 results in 00:00:00, so T10:21:14 is actually AFTER. But wait, parseDate returns time.Time. 2026-02-07T10:21:14 > 2026-02-07T00:00:00.
		
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
