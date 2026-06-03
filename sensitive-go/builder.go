package gosensitive

import (
	"github.com/Karrecy/sensitive-go/algorithm"
	"github.com/Karrecy/sensitive-go/algorithm/ac"
	"github.com/Karrecy/sensitive-go/algorithm/dfa"
	"github.com/Karrecy/sensitive-go/builtin"
	"github.com/Karrecy/sensitive-go/dict"
	"github.com/Karrecy/sensitive-go/filter"
	"github.com/Karrecy/sensitive-go/loader"
	"github.com/Karrecy/sensitive-go/variant"
)

// Builder provides a fluent API for constructing a Detector
type Builder struct {
	algorithmType    AlgorithmType
	words            []dict.Word
	loaders          []loader.Loader
	options          *Options
	whitelist        []string
	whitelistLoaders []loader.Loader      // Loaders for whitelist
	fileLoaders      []*loader.FileLoader // Track file loaders for watching
}

// New creates a new Builder with default settings
func New() *Builder {
	return &Builder{
		algorithmType:    AlgorithmAuto,
		words:            make([]dict.Word, 0),
		loaders:          make([]loader.Loader, 0),
		options:          DefaultOptions(),
		whitelist:        make([]string, 0),
		whitelistLoaders: make([]loader.Loader, 0),
		fileLoaders:      make([]*loader.FileLoader, 0),
	}
}

// UseAlgorithm specifies which matching algorithm to use
func (b *Builder) UseAlgorithm(algo AlgorithmType) *Builder {
	b.algorithmType = algo
	b.options.Algorithm = algo
	return b
}

// LoadFile adds a file loader to load words from the specified file path
func (b *Builder) LoadFile(path string) *Builder {
	fileLoader := loader.NewFileLoader(path)
	b.loaders = append(b.loaders, fileLoader)
	b.fileLoaders = append(b.fileLoaders, fileLoader)
	return b
}

// LoadMemory adds words directly from a string slice
func (b *Builder) LoadMemory(words []string) *Builder {
	b.loaders = append(b.loaders, loader.NewMemoryLoader(words))
	return b
}

// LoadHTTP adds an HTTP loader to load words from a remote URL
func (b *Builder) LoadHTTP(url string) *Builder {
	b.loaders = append(b.loaders, loader.NewHTTPLoader(url))
	return b
}

// LoadWords adds words directly to the builder
func (b *Builder) LoadWords(words []dict.Word) *Builder {
	b.words = append(b.words, words...)
	return b
}

// LoadBuiltin loads the built-in default word dictionary
func (b *Builder) LoadBuiltin() *Builder {
	b.words = append(b.words, builtin.GetDefaultWords()...)
	return b
}

// EnablePinyin enables pinyin variant detection
func (b *Builder) EnablePinyin() *Builder {
	b.options.EnablePinyin = true
	return b
}

// EnableVariant enables traditional Chinese variant detection
func (b *Builder) EnableVariant() *Builder {
	b.options.EnableTraditional = true
	return b
}

// EnableSymbol enables symbol interference filtering
func (b *Builder) EnableSymbol() *Builder {
	b.options.EnableSymbolFilter = true
	return b
}

// EnableSimilarChar enables similar character detection
func (b *Builder) EnableSimilarChar() *Builder {
	b.options.EnableSimilarChar = true
	return b
}

// AddWhitelist adds words to the whitelist
func (b *Builder) AddWhitelist(words ...string) *Builder {
	b.whitelist = append(b.whitelist, words...)
	return b
}

// LoadWhitelistFile loads whitelist from a file
func (b *Builder) LoadWhitelistFile(path string) *Builder {
	b.whitelistLoaders = append(b.whitelistLoaders, loader.NewFileLoader(path))
	return b
}

// LoadWhitelistMemory loads whitelist from memory
func (b *Builder) LoadWhitelistMemory(words []string) *Builder {
	b.whitelistLoaders = append(b.whitelistLoaders, loader.NewMemoryLoader(words))
	return b
}

// LoadWhitelistHTTP loads whitelist from a remote HTTP(S) URL
func (b *Builder) LoadWhitelistHTTP(url string) *Builder {
	b.whitelistLoaders = append(b.whitelistLoaders, loader.NewHTTPLoader(url))
	return b
}

// SetOptions sets custom options
func (b *Builder) SetOptions(opts *Options) *Builder {
	if opts != nil {
		b.options = opts
		b.algorithmType = opts.Algorithm
	}
	return b
}

// SetReplaceChar sets the replacement character
func (b *Builder) SetReplaceChar(r rune) *Builder {
	b.options.ReplaceChar = r
	return b
}

// SetCaseSensitive sets whether matching should be case-sensitive
func (b *Builder) SetCaseSensitive(sensitive bool) *Builder {
	b.options.CaseSensitive = sensitive
	return b
}

// Build constructs the Detector from the configured settings
func (b *Builder) Build() (*Detector, error) {
	// Load words from all loaders
	for _, l := range b.loaders {
		loadedWords, err := l.Load()
		if err != nil {
			return nil, err
		}
		b.words = append(b.words, loadedWords...)
	}

	// Choose algorithm based on word count if auto
	var matcher algorithm.Matcher
	if b.algorithmType == AlgorithmAuto {
		if len(b.words) < 5000 {
			matcher = dfa.NewDFAMatcher(b.options.CaseSensitive)
		} else {
			matcher = ac.NewACMatcher(b.options.CaseSensitive)
		}
	} else if b.algorithmType == AlgorithmDFA {
		matcher = dfa.NewDFAMatcher(b.options.CaseSensitive)
	} else {
		matcher = ac.NewACMatcher(b.options.CaseSensitive)
	}

	// Build the matcher
	if err := matcher.Build(b.words); err != nil {
		return nil, err
	}

	// Create detector
	detector := &Detector{
		matcher:    matcher,
		options:    b.options,
		filters:    make([]filter.Filter, 0),
		processors: make([]variant.Processor, 0),
		watchers:   make([]*FileWatcher, 0),
	}

	// Initialize variant processors based on options
	if b.options.EnableSymbolFilter {
		detector.processors = append(detector.processors, variant.NewSymbolProcessor())
	}
	if b.options.EnableTraditional {
		detector.processors = append(detector.processors, variant.NewTraditionalProcessor())
	}
	if b.options.EnableSimilarChar {
		detector.processors = append(detector.processors, variant.NewSimilarProcessor())
	}
	if b.options.EnablePinyin {
		detector.processors = append(detector.processors, variant.NewPinyinProcessor())
	}

	// Load whitelist from loaders
	whitelistWords := make([]string, 0)
	for _, l := range b.whitelistLoaders {
		loadedWords, err := l.Load()
		if err != nil {
			// Skip failed loaders but continue
			continue
		}
		// Extract text from Word objects
		for _, w := range loadedWords {
			whitelistWords = append(whitelistWords, w.Text)
		}
	}

	// Combine with directly added whitelist
	whitelistWords = append(whitelistWords, b.whitelist...)

	// Add whitelist filter if provided
	if len(whitelistWords) > 0 {
		detector.filters = append(detector.filters, filter.NewWhitelist(whitelistWords))
	}

	// Start file watchers if enabled
	if b.options.WatchFile && len(b.fileLoaders) > 0 {
		for _, fileLoader := range b.fileLoaders {
			watcher := NewFileWatcher(detector, fileLoader, b.options.WatchInterval)
			watcher.Start()
			detector.watchers = append(detector.watchers, watcher)
		}
	}

	return detector, nil
}
