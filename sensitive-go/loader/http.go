package loader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Karrecy/sensitive-go/dict"
)

// HTTPLoader loads sensitive words from a remote HTTP(S) URL
type HTTPLoader struct {
	url     string
	timeout time.Duration
	client  *http.Client
}

// NewHTTPLoader creates a new HTTP loader
func NewHTTPLoader(url string) *HTTPLoader {
	return &HTTPLoader{
		url:     url,
		timeout: 30 * time.Second,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetTimeout sets the HTTP request timeout
func (l *HTTPLoader) SetTimeout(timeout time.Duration) *HTTPLoader {
	l.timeout = timeout
	l.client.Timeout = timeout
	return l
}

// Load downloads and loads words from the URL
func (l *HTTPLoader) Load() ([]dict.Word, error) {
	resp, err := l.client.Get(l.url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	// Try to determine format from Content-Type
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		return l.loadJSON(resp.Body)
	}

	// Default to plain text
	return l.loadTXT(resp.Body)
}

// loadTXT loads words from plain text response
func (l *HTTPLoader) loadTXT(reader io.Reader) ([]dict.Word, error) {
	words := make([]dict.Word, 0)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		words = append(words, dict.Word{
			Text:     line,
			Category: dict.CategoryOther,
			Level:    dict.LevelMedium,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return words, nil
}

// loadJSON loads words from JSON response
func (l *HTTPLoader) loadJSON(reader io.Reader) ([]dict.Word, error) {
	var words []dict.Word
	decoder := json.NewDecoder(reader)

	if err := decoder.Decode(&words); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return words, nil
}


