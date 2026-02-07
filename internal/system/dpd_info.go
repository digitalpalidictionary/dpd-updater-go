package system

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DPDInfo struct {
	Path     string
	Date     time.Time
	Bookname string
}

func ParseIFO(path string) (*DPDInfo, error) {
	filename := strings.ToLower(filepath.Base(path))
	if !strings.HasPrefix(filename, "dpd") {
		return nil, fmt.Errorf("not a DPD file: %s", filename)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info := &DPDInfo{Path: path}
	scanner := bufio.NewScanner(file)
	var dateStr string
	isDPD := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "bookname=") {
			info.Bookname = strings.TrimPrefix(line, "bookname=")
		} else if strings.HasPrefix(line, "date=") {
			dateStr = strings.TrimPrefix(line, "date=")
		} else if strings.HasPrefix(line, "author=") {
			author := strings.ToLower(strings.TrimPrefix(line, "author="))
			if strings.Contains(author, "bodhirasa") {
				isDPD = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Double check to avoid picking up unrelated dictionaries
	lowerBName := strings.ToLower(info.Bookname)
	if !isDPD && !strings.Contains(lowerBName, "pāḷi dictionary") && !strings.Contains(lowerBName, "dpd") {
		return nil, fmt.Errorf("not a verified DPD dictionary: %s", path)
	}

	if dateStr == "" {
		return nil, fmt.Errorf("date field not found in %s", path)
	}

	// Try parsing with time only first, if that fails try standard formats
	// The example showed "2026-02-06T14:40:12" which is ISO8601-like
	// Adjust format based on actual file content if needed
	parsedDate, err := time.Parse("2006-01-02T15:04:05", dateStr)
	if err != nil {
		// Fallback or try other formats if needed
		return nil, fmt.Errorf("failed to parse date '%s': %w", dateStr, err)
	}
	info.Date = parsedDate

	return info, nil
}
