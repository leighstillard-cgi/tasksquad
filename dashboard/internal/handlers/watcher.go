package handlers

import (
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher monitors worklog files and triggers refresh on changes.
type Watcher struct {
	watcher      *fsnotify.Watcher
	worklogPath  string
	logger       *slog.Logger
	onChangeFunc func()
	debounceTime time.Duration

	debounceTimer *time.Timer
	debounceMu    sync.Mutex
}

// NewWatcher creates a file watcher for the worklog directory.
func NewWatcher(worklogPath string, logger *slog.Logger, onChangeFunc func()) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		watcher:      fsWatcher,
		worklogPath:  worklogPath,
		logger:       logger,
		onChangeFunc: onChangeFunc,
		debounceTime: 100 * time.Millisecond,
	}

	return w, nil
}

// Start begins watching the configured paths.
func (w *Watcher) Start() error {
	// Watch individual files
	filesToWatch := []string{
		filepath.Join(w.worklogPath, "dispatch-log.md"),
		filepath.Join(w.worklogPath, "backlog.md"),
	}

	// Watch directories
	dirsToWatch := []string{
		filepath.Join(w.worklogPath, "completions"),
		filepath.Join(w.worklogPath, "escalations"),
		filepath.Join(w.worklogPath, "dispatches"),
		filepath.Join(w.worklogPath, "session-logs"),
	}

	for _, f := range filesToWatch {
		if err := w.watcher.Add(f); err != nil {
			w.logger.Warn("could not watch file", "path", f, "error", err)
		} else {
			w.logger.Debug("watching file", "path", f)
		}
	}

	for _, d := range dirsToWatch {
		if err := w.watcher.Add(d); err != nil {
			w.logger.Warn("could not watch directory", "path", d, "error", err)
		} else {
			w.logger.Debug("watching directory", "path", d)
		}
	}

	go w.eventLoop()
	return nil
}

// Stop stops the watcher.
func (w *Watcher) Stop() error {
	return w.watcher.Close()
}

func (w *Watcher) eventLoop() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			w.handleEvent(event)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			w.logger.Error("watcher error", "error", err)
		}
	}
}

func (w *Watcher) handleEvent(event fsnotify.Event) {
	// Only react to write/create/remove events
	if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) == 0 {
		return
	}

	w.logger.Debug("file change detected", "path", event.Name, "op", event.Op.String())
	w.debounce()
}

func (w *Watcher) debounce() {
	w.debounceMu.Lock()
	defer w.debounceMu.Unlock()

	if w.debounceTimer != nil {
		w.debounceTimer.Stop()
	}

	w.debounceTimer = time.AfterFunc(w.debounceTime, func() {
		w.logger.Info("triggering refresh from file change")
		w.onChangeFunc()
	})
}
