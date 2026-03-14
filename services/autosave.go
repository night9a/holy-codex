package services

import (
	"log"
	"sync"
	"time"

	"holy-codex/domain"
	"holy-codex/infrastructure/storage"
)

// AutoSave buffers pending entries and flushes them on a timer.
type AutoSave struct {
	store    storage.Storage
	interval time.Duration

	mu      sync.Mutex
	pending map[string]*domain.DiaryEntry

	ticker *time.Ticker
	stopCh chan struct{}
}

// NewAutoSave constructs an AutoSave with a 3-second default flush interval.
func NewAutoSave(store storage.Storage) *AutoSave {
	return &AutoSave{
		store:    store,
		interval: 3 * time.Second,
		pending:  make(map[string]*domain.DiaryEntry),
		stopCh:   make(chan struct{}),
	}
}

// Stage queues an entry to be persisted on the next tick.
func (a *AutoSave) Stage(entry *domain.DiaryEntry) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.pending[entry.ID] = entry
}

// Flush immediately persists all staged entries.
func (a *AutoSave) Flush() {
	a.mu.Lock()
	snapshot := make(map[string]*domain.DiaryEntry, len(a.pending))
	for k, v := range a.pending {
		snapshot[k] = v
	}
	a.pending = make(map[string]*domain.DiaryEntry)
	a.mu.Unlock()

	for _, e := range snapshot {
		if err := a.store.SaveEntry(e); err != nil {
			log.Printf("[autosave] failed to save %s: %v", e.ID, err)
		}
	}
}

// Start launches the background flush goroutine.
func (a *AutoSave) Start() {
	a.ticker = time.NewTicker(a.interval)
	go func() {
		for {
			select {
			case <-a.stopCh:
				return
			case <-a.ticker.C:
				a.Flush()
			}
		}
	}()
}

// Stop halts auto-save and does a final flush.
func (a *AutoSave) Stop() {
	if a.ticker != nil {
		a.ticker.Stop()
	}
	close(a.stopCh)
	a.Flush()
}