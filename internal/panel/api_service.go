// Package panel exports Wails-bound API methods for the frontend.
package panel

import (
	"context"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

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

// log emits a log line event to the frontend.
func (a *API) log(event, line string) {
	if a.ctx == nil {
		return
	}
	runtime.EventsEmit(a.ctx, event, line)
}

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

// TestRoom tests room connectivity. Emits "test:log" events with progress.
// Returns "ok" or error message.
func (a *API) TestRoom(provider, transport, roomID string) string {
	a.log("test:log", "testing "+provider+"/"+transport+" room: "+roomID)
	a.log("test:log", "connecting to room...")

	res, err := rooms.TestRoom(a.ctx, provider, transport, roomID)
	if err != nil {
		a.log("test:log", "error: "+err.Error())
		return err.Error()
	}
	if res == nil {
		a.log("test:log", "error: internal error")
		return "internal error"
	}
	if res.OK {
		a.log("test:log", "ping/pong ok")
		return "ok"
	}
	a.log("test:log", "failed: "+res.Message)
	return res.Message
}

// CreateInstance creates and persists a new olcrtc instance.
func (a *API) CreateInstance(inst instances.Instance) instances.Instance {
	if inst.Key == "" {
		inst.Key = instances.NewKey()
	}
	if err := instances.Add(&inst); err != nil {
		a.log("create:error", err.Error())
	}
	return inst
}

// ListInstances returns all persisted instances.
func (a *API) ListInstances() []instances.Instance {
	list, err := instances.LoadAll()
	if err != nil {
		return []instances.Instance{}
	}
	return list
}

// DeleteInstance removes an instance by ID.
func (a *API) DeleteInstance(id string) string {
	if err := instances.Remove(id); err != nil {
		return err.Error()
	}
	return "ok"
}

// Install installs olcrtc on the server via SSH. Emits "install:log" events.
func (a *API) Install(host string, port int, user, password string) string {
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Minute)
	defer cancel()

	a.log("install:log", "connecting to "+user+"@"+host+":"+fmt.Sprint(port))

	runner, err := exec.NewSSHRunner(exec.SSHConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	})
	if err != nil {
		msg := "ssh: " + err.Error()
		a.log("install:log", msg)
		return msg
	}

	a.log("install:log", "connected, starting install...")

	if err := installer.New(runner).WithLog(func(line string) {
		a.log("install:log", line)
	}).Install(ctx); err != nil {
		msg := fmt.Sprintf("install failed: %v", err)
		a.log("install:log", msg)
		return msg
	}

	a.log("install:log", "installation complete")
	a.log("install:done", "ok")
	return "ok"
}
