package panel

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"src.owenewans.org/owenrtc/internal/instances"
	"src.owenewans.org/owenrtc/internal/rooms"
)

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	p := filepath.Join(s.web, r.URL.Path)
	if r.URL.Path == "/" || r.URL.Path == "" {
		p = filepath.Join(s.web, "index.html")
	}
	http.ServeFile(w, r, p)
}

func (s *Server) handleMode(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]any{
		"mode":      string(s.mode.Kind),
		"port":      s.mode.Port,
		"public_ip": s.mode.PublicIP,
	})
}

func (s *Server) handleServers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, []instances.Instance{})
	case http.MethodPost:
		var inst instances.Instance
		if err := json.NewDecoder(r.Body).Decode(&inst); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		writeJSON(w, inst)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleJitsiRooms(w http.ResponseWriter, r *http.Request) {
	hosts, err := rooms.LoadJitsiInstances("olcrtc")
	if err != nil {
		http.Error(w, "failed to load jitsi instances", http.StatusInternalServerError)
		return
	}
	writeJSON(w, hosts)
}

func (s *Server) handleTestRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Provider  string `json:"provider"`
		Transport string `json:"transport"`
		RoomID    string `json:"room_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	result, err := rooms.TestRoom(r.Context(), req.Provider, req.Transport, req.RoomID)
	if err != nil {
		writeJSON(w, map[string]string{"result": err.Error()})
		return
	}
	msg := result.Message
	if result.OK {
		msg = "ok"
	}
	writeJSON(w, map[string]string{"result": msg})
}

func (s *Server) handleInstall(w http.ResponseWriter, r *http.Request) {
	// TODO: wire to installer with correct runner (self or ssh)
	writeJSON(w, map[string]string{"status": "not implemented"})
}
