package gosensitive

import (
	"sync"

	"github.com/Karrecy/sensitive-go/algorithm"
	"github.com/Karrecy/sensitive-go/algorithm/ac"
	"github.com/Karrecy/sensitive-go/algorithm/dfa"
	"github.com/Karrecy/sensitive-go/dict"
	"github.com/Karrecy/sensitive-go/filter"
	"github.com/Karrecy/sensitive-go/variant"
)

// Detector is the main sensitive word detector
type Detector struct {
	matcher    algorithm.Matcher
	filters    []filter.Filter
	processors []variant.Processor
	watchers   []*FileWatcher
	options    *Options
	mu         sync.RWMutex
}

// preprocessText applies all enabled variant processors to normalize text
func (d *Detector) preprocessText(text string) string {
	for _, processor := range d.processors {
		text = processor.Process(text)
	}
	return text
}

// Contains checks if the text contains any sensitive words
func (d *Detector) Contains(text string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Preprocess text with variant processors
	text = d.preprocessText(text)
	
	return !d.matcher.Validate(text)
}

// Find returns all sensitive words found in the text
func (d *Detector) Find(text string) []Match {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Preprocess text with variant processors
	text = d.preprocessText(text)
	
	matches := d.matcher.Match(text)
	result := make([]Match, 0, len(matches))

	for _, m := range matches {
		// Apply filters
		if d.shouldFilter(m.Word) {
			continue
		}

		// Apply category filter
		if len(d.options.Categories) > 0 {
			found := false
			for _, cat := range d.options.Categories {
				if m.Category&dict.Category(cat) != 0 {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Apply level filter
		if m.Level < dict.Level(d.options.MinLevel) {
			continue
		}

		result = append(result, Match{
			Word:     m.Word,
			Start:    m.Start,
			End:      m.End,
			Category: dict.Category(m.Category),
			Level:    dict.Level(m.Level),
		})

		// Check max match count
		if d.options.MaxMatchCount > 0 && len(result) >= d.options.MaxMatchCount {
			break
		}
	}

	return result
}

// FindAll returns detailed detection results
func (d *Detector) FindAll(text string) *Result {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Preprocess text once
	processedText := d.preprocessText(text)
	
	// Match on preprocessed text
	matches := d.matcher.Match(processedText)
	result := make([]Match, 0, len(matches))

	for _, m := range matches {
		// Apply filters
		if d.shouldFilter(m.Word) {
			continue
		}

		// Apply category filter
		if len(d.options.Categories) > 0 {
			found := false
			for _, cat := range d.options.Categories {
				if m.Category&dict.Category(cat) != 0 {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Apply level filter
		if m.Level < dict.Level(d.options.MinLevel) {
			continue
		}

		result = append(result, Match{
			Word:     m.Word,
			Start:    m.Start,
			End:      m.End,
			Category: dict.Category(m.Category),
			Level:    dict.Level(m.Level),
		})

		// Check max match count
		if d.options.MaxMatchCount > 0 && len(result) >= d.options.MaxMatchCount {
			break
		}
	}
	
	// Get filtered text
	filteredText := d.matcher.Replace(processedText, d.options.ReplaceChar)
	
	return &Result{
		Found:        len(result) > 0,
		Matches:      result,
		FilteredText: filteredText,
	}
}

// Replace replaces sensitive words with the given replacement string
func (d *Detector) Replace(text, replacement string) string {
	if len(replacement) == 0 {
		return text
	}
	
	// Use the first rune of replacement
	repl := []rune(replacement)[0]
	return d.ReplaceRune(text, repl)
}

// ReplaceRune replaces sensitive words with the given rune
func (d *Detector) ReplaceRune(text string, repl rune) string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Preprocess text with variant processors
	text = d.preprocessText(text)
	
	return d.matcher.Replace(text, repl)
}

// Validate checks if the text is clean (returns true if no sensitive words found)
func (d *Detector) Validate(text string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Preprocess text with variant processors
	text = d.preprocessText(text)
	
	return d.matcher.Validate(text)
}

// Filter returns the text with sensitive words replaced
func (d *Detector) Filter(text string) string {
	return d.ReplaceRune(text, d.options.ReplaceChar)
}

// Reload reloads the detector with new words atomically
// If reload fails, the detector keeps using the old matcher
func (d *Detector) Reload(words []dict.Word) error {
	// Create a new matcher instance
	var newMatcher algorithm.Matcher
	caseSensitive := d.options.CaseSensitive
	
	// Determine algorithm type based on word count
	if len(words) < 5000 {
		newMatcher = dfa.NewDFAMatcher(caseSensitive)
	} else {
		newMatcher = ac.NewACMatcher(caseSensitive)
	}
	
	// Build the new matcher
	if err := newMatcher.Build(words); err != nil {
		return err
	}
	
	// Atomically replace the old matcher with the new one
	d.mu.Lock()
	d.matcher = newMatcher
	d.mu.Unlock()
	
	return nil
}

// shouldFilter checks if a word should be filtered out by whitelist
func (d *Detector) shouldFilter(word string) bool {
	for _, f := range d.filters {
		if f.ShouldFilter(word) {
			return true
		}
	}
	return false
}

// AddFilter adds a custom filter
func (d *Detector) AddFilter(f filter.Filter) {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.filters = append(d.filters, f)
}

// Close stops all file watchers and releases resources
func (d *Detector) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	for _, watcher := range d.watchers {
		watcher.Stop()
	}
	d.watchers = nil
	
	return nil
}

