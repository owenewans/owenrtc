// Package installer downloads olcrtc binaries from owenrtc gh releases.
//
// Install never builds olcrtc locally - it fetches pre-built binaries
// from https://github.com/owenewans/owenrtc/releases.
package installer

import (
	"context"
	"fmt"
	"runtime"

	"src.owenewans.org/owenrtc/internal/exec"
)

// LogFunc is called for each install step.
type LogFunc func(line string)

// Installer installs olcrtc on a server (local or remote via SSH).
type Installer struct {
	runner exec.Runner
	log    LogFunc
}

// New creates an installer with the given command runner.
func New(r exec.Runner) *Installer {
	return &Installer{runner: r, log: func(string) {}}
}

// WithLog sets a log callback for install progress.
func (i *Installer) WithLog(fn LogFunc) *Installer {
	i.log = fn
	return i
}

// Install downloads and installs the latest olcrtc binary.
func (i *Installer) Install(ctx context.Context) error {
	arch := runtime.GOARCH
	if runtime.GOOS != "linux" {
		return fmt.Errorf("install: only linux supported, got %s", runtime.GOOS)
	}
	url := fmt.Sprintf(
		"https://github.com/owenewans/owenrtc/releases/latest/download/olcrtc-linux-%s",
		arch,
	)
	steps := []struct {
		desc string
		cmd  string
	}{
		{"creating directory", "mkdir -p ~/.owenrtc/bin"},
		{"downloading binary", fmt.Sprintf("curl -sL %s -o ~/.owenrtc/bin/olcrtc", url)},
		{"setting permissions", "chmod +x ~/.owenrtc/bin/olcrtc"},
		{"verifying install", "~/.owenrtc/bin/olcrtc --help 2>&1 | head -1 || true"},
		{"adding to PATH", "grep -q 'owenrtc/bin' ~/.bashrc 2>/dev/null || echo 'export PATH=$HOME/.owenrtc/bin:$PATH' >> ~/.bashrc"},
	}
	for _, s := range steps {
		i.log(s.desc + "...")
		if _, err := i.runner.Run(ctx, s.cmd); err != nil {
			return fmt.Errorf("step %q: %w", s.desc, err)
		}
	}
	i.log("install done")
	return nil
}
