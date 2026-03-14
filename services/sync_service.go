package services

import (
	"log"
	"time"

	"holy-codex/infrastructure/network"
)

// SyncService periodically triggers a push/pull cycle with peers.
type SyncService struct {
	syncer   *network.Sync
	interval time.Duration
	ticker   *time.Ticker
	stopCh   chan struct{}
}

// NewSyncService creates a SyncService that runs every 30 seconds.
func NewSyncService(syncer *network.Sync) *SyncService {
	return &SyncService{
		syncer:   syncer,
		interval: 30 * time.Second,
		stopCh:   make(chan struct{}),
	}
}

// Start begins the sync loop and starts the HTTP listener on the syncer.
func (s *SyncService) Start() {
	s.syncer.StartServer()
	s.ticker = time.NewTicker(s.interval)
	go func() {
		for {
			select {
			case <-s.stopCh:
				return
			case <-s.ticker.C:
				if err := s.syncer.PushToAll(); err != nil {
					log.Printf("[sync_service] push error: %v", err)
				}
			}
		}
	}()
}

// Stop halts the sync loop and shuts down the HTTP server.
func (s *SyncService) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	close(s.stopCh)
	s.syncer.StopServer()
}