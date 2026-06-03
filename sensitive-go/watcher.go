package gosensitive

import (
	"os"
	"sync"
	"time"

	"github.com/Karrecy/sensitive-go/loader"
)

// FileWatcher monitors file changes and triggers reload
type FileWatcher struct {
	detector *Detector
	loader   loader.Loader
	interval time.Duration
	lastMod  time.Time
	stopCh   chan struct{}
	mu       sync.Mutex
	running  bool
}

// NewFileWatcher creates a new file watcher
func NewFileWatcher(detector *Detector, loader loader.Loader, interval time.Duration) *FileWatcher {
	return &FileWatcher{
		detector: detector,
		loader:   loader,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start begins monitoring the file for changes
func (w *FileWatcher) Start() error {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return nil
	}
	w.running = true
	w.mu.Unlock()

	// Get initial file modification time
	if fileLoader, ok := w.loader.(*loader.FileLoader); ok {
		if info, err := os.Stat(fileLoader.Path()); err == nil {
			w.lastMod = info.ModTime()
		}
	}

	go w.watch()
	return nil
}

// Stop stops the file watcher
func (w *FileWatcher) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.running {
		return
	}

	w.running = false
	close(w.stopCh)
}

// watch monitors file changes
func (w *FileWatcher) watch() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.checkAndReload()
		case <-w.stopCh:
			return
		}
	}
}

// checkAndReload checks if file has changed and reloads if necessary
func (w *FileWatcher) checkAndReload() {
	fileLoader, ok := w.loader.(*loader.FileLoader)
	if !ok {
		return
	}

	info, err := os.Stat(fileLoader.Path())
	if err != nil {
		// File doesn't exist or can't be accessed
		return
	}

	modTime := info.ModTime()
	if modTime.After(w.lastMod) {
		// File has been modified, reload
		w.lastMod = modTime
		w.reloadWords()
	}
}

// reloadWords reloads words from the file
func (w *FileWatcher) reloadWords() {
	words, err := w.loader.Load()
	if err != nil {
		// Failed to load, keep using current words
		return
	}

	// Reload the detector with new words
	w.detector.Reload(words)
}

// IsRunning returns whether the watcher is currently running
func (w *FileWatcher) IsRunning() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.running
}
