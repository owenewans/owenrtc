// Package panel serves the owenrtc web UI and API.
package panel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"src.owenewans.org/owenrtc/internal/mode"
)

// Server is the owenrtc panel HTTP server.
type Server struct {
	mode mode.Mode
	mux  *http.ServeMux
	web  string
}

// New creates a panel server.
func New(m mode.Mode) *Server {
	s := &Server{
		mode: m,
		mux:  http.NewServeMux(),
		web:  webDir(),
	}
	s.routes()
	return s
}

func webDir() string {
	if _, err := os.Stat("web"); err == nil {
		return "web"
	}
	exe, err := os.Executable()
	if err != nil {
		return "web"
	}
	return filepath.Join(filepath.Dir(exe), "web")
}

func (s *Server) routes() {
	s.mux.HandleFunc("/", s.handleStatic)
	s.mux.HandleFunc("/api/mode", s.handleMode)
	s.mux.HandleFunc("/api/servers", s.handleServers)
	s.mux.HandleFunc("/api/rooms/jitsi", s.handleJitsiRooms)
	s.mux.HandleFunc("/api/rooms/test", s.handleTestRoom)
	s.mux.HandleFunc("/api/install", s.handleInstall)
}

// Run starts the HTTP server.
func (s *Server) Run(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.mode.Port)
	srv := &http.Server{Addr: addr, Handler: s.mux}

	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()

	log.Printf("panel on http://127.0.0.1:%d", s.mode.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("panel: %w", err)
	}
	return nil
}

// Mux returns the underlying mux so external callers (e.g. wails) can route API calls.
func (s *Server) Mux() *http.ServeMux {
	return s.mux
}

// IsAPI reports whether path targets the API.
func IsAPI(p string) bool {
	return strings.HasPrefix(p, "/api/")
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
