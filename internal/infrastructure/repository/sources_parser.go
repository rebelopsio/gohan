package repository

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	domainRepo "github.com/rebelopsio/gohan/internal/domain/repository"
)

var (
	// ErrInvalidLine is returned when a line cannot be parsed
	ErrInvalidLine = errors.New("invalid sources.list line")
)

// ParsedEntry represents a parsed sources.list entry
type ParsedEntry = domainRepo.SourceEntry

// ParseSourcesFile parses the content of a sources.list file
func ParseSourcesFile(content string) ([]ParsedEntry, error) {
	var entries []ParsedEntry

	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		entry, err := ParseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}

		// Skip nil entries (comments, empty lines)
		if entry != nil {
			entries = append(entries, *entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read content: %w", err)
	}

	return entries, nil
}

// ParseLine parses a single line from sources.list
// Returns nil for comments and empty lines
func ParseLine(line string) (*ParsedEntry, error) {
	// Trim whitespace
	line = strings.TrimSpace(line)

	// Skip empty lines
	if line == "" {
		return nil, nil
	}

	// Skip comments
	if strings.HasPrefix(line, "#") {
		return nil, nil
	}

	// Split line into fields
	fields := strings.Fields(line)

	// Minimum: type URI suite component1 [component2...]
	if len(fields) < 4 {
		return nil, fmt.Errorf("%w: expected at least 4 fields, got %d", ErrInvalidLine, len(fields))
	}

	entryType := fields[0]
	uri := fields[1]
	suite := fields[2]
	components := fields[3:]

	// Validate type
	if entryType != "deb" && entryType != "deb-src" {
		return nil, fmt.Errorf("%w: invalid type '%s'", ErrInvalidLine, entryType)
	}

	return &ParsedEntry{
		Type:       entryType,
		URI:        uri,
		Suite:      suite,
		Components: components,
	}, nil
}
