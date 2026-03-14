package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"bytes"

	"holy-codex/domain"
	"holy-codex/infrastructure/storage"
)

// SyncPayload is what peers exchange over HTTP.
type SyncPayload struct {
	Entries []*domain.DiaryEntry `json:"entries"`
}

// Sync pushes unsynced local entries to all known peers and pulls theirs.
type Sync struct {
	store     storage.Storage
	discovery *Discovery
	port      int
	server    *http.Server
}

// NewSync constructs a Sync instance.
func NewSync(store storage.Storage, discovery *Discovery, port int) *Sync {
	return &Sync{
		store:     store,
		discovery: discovery,
		port:      port,
	}
}

// StartServer starts the HTTP listener that peers push entries to.
func (s *Sync) StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/sync/push", s.handlePush)
	mux.HandleFunc("/sync/pull", s.handlePull)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[sync] server error: %v", err)
		}
	}()
	log.Printf("[sync] listening on :%d", s.port)
}

// StopServer gracefully shuts down the HTTP sync server.
func (s *Sync) StopServer() {
	if s.server != nil {
		_ = s.server.Close()
	}
}

// PushToAll sends local unsynced entries to every discovered peer.
func (s *Sync) PushToAll() error {
	entries, err := s.store.UnsyncedEntries()
	if err != nil || len(entries) == 0 {
		return err
	}

	payload, err := json.Marshal(SyncPayload{Entries: entries})
	if err != nil {
		return err
	}

	for _, peer := range s.discovery.Peers() {
		url := fmt.Sprintf("http://%s:%d/sync/push", peer.Addr, peer.Port)
		resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
		if err != nil {
			log.Printf("[sync] push to %s failed: %v", peer.Name, err)
			continue
		}
		resp.Body.Close()
	}

	for _, e := range entries {
		_ = s.store.MarkSynced(e.ID)
	}
	return nil
}

// ─── HTTP handlers ────────────────────────────────────────────────────────────

func (s *Sync) handlePush(w http.ResponseWriter, r *http.Request) {
	var payload SyncPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad payload", http.StatusBadRequest)
		return
	}
	for _, e := range payload.Entries {
		if err := s.store.SaveEntry(e); err != nil {
			log.Printf("[sync] failed to store entry %s: %v", e.ID, err)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Sync) handlePull(w http.ResponseWriter, r *http.Request) {
	entries, err := s.store.UnsyncedEntries()
	if err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(SyncPayload{Entries: entries})
}

// ─── Helpers ──────────────────────────────────────────────────────────────────


func bytesReader(b []byte) *bytesBuf { return &bytesBuf{b: b} }

type bytesBuf struct {
	b   []byte
	pos int
}

func (bb *bytesBuf) Read(p []byte) (int, error) {
	if bb.pos >= len(bb.b) {
		return 0, fmt.Errorf("EOF")
	}
	n := copy(p, bb.b[bb.pos:])
	bb.pos += n
	return n, nil
}