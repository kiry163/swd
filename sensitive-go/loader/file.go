package loader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Karrecy/sensitive-go/dict"
)

// FileLoader loads sensitive words from a file
type FileLoader struct {
	path string
}

// NewFileLoader creates a new file loader
func NewFileLoader(path string) *FileLoader {
	return &FileLoader{path: path}
}

// Load loads words from the file
func (l *FileLoader) Load() ([]dict.Word, error) {
	ext := filepath.Ext(l.path)
	
	switch strings.ToLower(ext) {
	case ".json":
		return l.loadJSON()
	case ".txt":
		return l.loadTXT()
	default:
		return l.loadTXT() // Default to txt format
	}
}

// loadTXT loads words from a plain text file (one word per line)
func (l *FileLoader) loadTXT() ([]dict.Word, error) {
	file, err := os.Open(l.path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	words := make([]dict.Word, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		words = append(words, dict.Word{
			Text:     line,
			Category: dict.CategoryOther,
			Level:    dict.LevelMedium,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return words, nil
}

// loadJSON loads words from a JSON file
func (l *FileLoader) loadJSON() ([]dict.Word, error) {
	file, err := os.Open(l.path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var words []dict.Word
	decoder := json.NewDecoder(file)
	
	if err := decoder.Decode(&words); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return words, nil
}

// Path returns the file path
func (l *FileLoader) Path() string {
	return l.path
}


