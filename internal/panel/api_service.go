// Package panel exports Wails-bound API methods for the frontend.
package panel

import (
	"context"
	"fmt"
	"time"

	"src.owenewans.org/owenrtc/internal/exec"
	"src.owenewans.org/owenrtc/internal/installer"
	"src.owenewans.org/owenrtc/internal/instances"
	"src.owenewans.org/owenrtc/internal/mode"
	"src.owenewans.org/owenrtc/internal/rooms"
)

// API is bound to the Wails frontend. Methods are called from JS via
// ../wailsjs/go/panel/API.js
type API struct {
	ctx context.Context
}

// NewAPI creates a Wails-bound API service.
func NewAPI() *API { return &API{} }

// Startup stores the Wails context.
func (a *API) Startup(ctx context.Context) { a.ctx = ctx }

// Mode returns the detected runtime mode.
func (a *API) Mode() mode.Mode {
	return mode.Detect()
}

// JitsiHosts returns the list of jitsi instances from olcrtc submodule.
func (a *API) JitsiHosts() []string {
	hosts, err := rooms.LoadJitsiInstances("olcrtc")
	if err != nil {
		return []string{}
	}
	return hosts
}

// TestRoom tests room connectivity via olcrtc ping/pong.
// Returns "ok" on success or the error message.
func (a *API) TestRoom(provider, transport, roomID string) string {
	res, err := rooms.TestRoom(a.ctx, provider, transport, roomID)
	if err != nil {
		return err.Error()
	}
	if res == nil {
		return "internal error"
	}
	if res.OK {
		return "ok"
	}
	return res.Message
}

// CreateInstance creates a new olcrtc instance.
func (a *API) CreateInstance(inst instances.Instance) instances.Instance {
	// TODO: persist + start
	return inst
}

// ListInstances returns all instances.
func (a *API) ListInstances() []instances.Instance {
	// TODO: load from config
	return []instances.Instance{}
}

// Install installs olcrtc on the server (local or remote) via SSH.
func (a *API) Install(host string, port int, user, password string) string {
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Minute)
	defer cancel()

	runner, err := exec.NewSSHRunner(exec.SSHConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	})
	if err != nil {
		return "ssh: " + err.Error()
	}

	if err := installer.New(runner).Install(ctx); err != nil {
		return fmt.Sprintf("install: %v", err)
	}
	return "ok"
}
