package network

import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
)

const (
	multicastAddr    = "224.0.0.251:5353" // mDNS-like multicast group
	beaconInterval   = 5 * time.Second
	peerTimeout      = 15 * time.Second
	discoveryService = "_holydiary._tcp"
)

// PeerInfo describes a discovered peer on the LAN.
type PeerInfo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Addr      string    `json:"addr"`
	Port      int       `json:"port"`
	SeenAt    time.Time `json:"-"`
}

// Discovery manages LAN peer discovery via UDP multicast beacons.
type Discovery struct {
	selfID   string
	selfName string
	port     int

	mu    sync.RWMutex
	peers map[string]*PeerInfo

	stopCh chan struct{}
}

// NewDiscovery creates a Discovery instance for the given sync port.
func NewDiscovery(port int) *Discovery {
	return &Discovery{
		selfID:   generateNodeID(),
		selfName: localHostname(),
		port:     port,
		peers:    make(map[string]*PeerInfo),
		stopCh:   make(chan struct{}),
	}
}

// Start begins broadcasting beacons and listening for peers.
func (d *Discovery) Start() {
	go d.broadcast()
	go d.listen()
	go d.reaper()
}

// Stop halts discovery.
func (d *Discovery) Stop() {
	close(d.stopCh)
}

// Peers returns a snapshot of currently known live peers.
func (d *Discovery) Peers() []*PeerInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()
	peers := make([]*PeerInfo, 0, len(d.peers))
	for _, p := range d.peers {
		peers = append(peers, p)
	}
	return peers
}

// ─── Internal ─────────────────────────────────────────────────────────────────

func (d *Discovery) broadcast() {
	ticker := time.NewTicker(beaconInterval)
	defer ticker.Stop()
	for {
		select {
		case <-d.stopCh:
			return
		case <-ticker.C:
			d.sendBeacon()
		}
	}
}

func (d *Discovery) sendBeacon() {
	addr, err := net.ResolveUDPAddr("udp4", multicastAddr)
	if err != nil {
		return
	}
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	payload, _ := json.Marshal(PeerInfo{
		ID:   d.selfID,
		Name: d.selfName,
		Port: d.port,
	})
	_, _ = conn.Write(payload)
}

func (d *Discovery) listen() {
	addr, err := net.ResolveUDPAddr("udp4", multicastAddr)
	if err != nil {
		log.Printf("[discovery] resolve error: %v", err)
		return
	}
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Printf("[discovery] listen error: %v", err)
		return
	}
	defer conn.Close()
	_ = conn.SetReadBuffer(1024 * 8)

	buf := make([]byte, 1024)
	for {
		select {
		case <-d.stopCh:
			return
		default:
		}
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		var p PeerInfo
		if err := json.Unmarshal(buf[:n], &p); err != nil {
			continue
		}
		if p.ID == d.selfID {
			continue // ignore own beacon
		}
		p.Addr = src.IP.String()
		p.SeenAt = time.Now()

		d.mu.Lock()
		d.peers[p.ID] = &p
		d.mu.Unlock()
	}
}

// reaper removes peers that have not been heard from within peerTimeout.
func (d *Discovery) reaper() {
	ticker := time.NewTicker(peerTimeout)
	defer ticker.Stop()
	for {
		select {
		case <-d.stopCh:
			return
		case <-ticker.C:
			cutoff := time.Now().Add(-peerTimeout)
			d.mu.Lock()
			for id, p := range d.peers {
				if p.SeenAt.Before(cutoff) {
					delete(d.peers, id)
				}
			}
			d.mu.Unlock()
		}
	}
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func generateNodeID() string {
	// TODO: persist to config so the same node keeps the same ID across restarts
	return randomHex(8)
}

func localHostname() string {
	h, err := net.LookupAddr("")
	if err != nil || len(h) == 0 {
		return "unknown-scribe"
	}
	return h[0]
}

func randomHex(n int) string {
	const hex = "0123456789abcdef"
	b := make([]byte, n)
	for i := range b {
		b[i] = hex[int(time.Now().UnixNano())%(len(hex))]
	}
	return string(b)
}