package telemetry

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type StatusResponse struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	UptimeSec int64     `json:"uptime_sec"`
	Peers     int       `json:"peers"`
	SyncState string    `json:"sync_state"`
	Timestamp time.Time `json:"timestamp"`
}

type Server struct {
	startTime time.Time
	mu        sync.RWMutex
	peers     int
}

func NewServer() *Server {
	return &Server{
		startTime: time.Now(),
		peers:     0,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/status", s.handleStatus)
	mux.HandleFunc("/metrics", s.handleMetrics)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	peers := s.peers
	s.mu.RUnlock()

	resp := StatusResponse{
		Status:    "active",
		Version:   "0.2.0",
		UptimeSec: int64(time.Since(s.startTime).Seconds()),
		Peers:     peers,
		SyncState: "synced",
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("# HELP kryvora_node_up Node operational status\n# TYPE kryvora_node_up gauge\nkryvora_node_up 1\n"))
}
