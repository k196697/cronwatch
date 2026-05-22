// Package healthcheck provides an HTTP endpoint for inspecting
// the current status of monitored cron jobs.
package healthcheck

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cronwatch/internal/metrics"
)

// Server exposes a lightweight HTTP health endpoint.
type Server struct {
	addr    string
	metrics *metrics.Store
	server  *http.Server
}

// StatusResponse is the JSON payload returned by the health endpoint.
type StatusResponse struct {
	Status string                    `json:"status"`
	Jobs   map[string]metrics.Snapshot `json:"jobs"`
	Time   time.Time                 `json:"time"`
}

// New creates a Server listening on addr.
func New(addr string, m *metrics.Store) *Server {
	s := &Server{
		addr:    addr,
		metrics: m,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/status", s.handleStatus)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return s
}

// Start begins serving HTTP requests. It blocks until the server stops.
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Stop gracefully shuts down the HTTP server.
func (s *Server) Stop() error {
	return s.server.Close()
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	all := s.metrics.All()
	resp := StatusResponse{
		Status: overallStatus(all),
		Jobs:   all,
		Time:   time.Now().UTC(),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func overallStatus(jobs map[string]metrics.Snapshot) string {
	for _, s := range jobs {
		if s.LastFailed {
			return "degraded"
		}
	}
	return "ok"
}
