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

// Installer installs olcrtc on a server (local or remote via SSH).
type Installer struct {
	runner exec.Runner
}

// New creates an installer with the given command runner.
func New(r exec.Runner) *Installer {
	return &Installer{runner: r}
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
	steps := []string{
		"mkdir -p ~/.owenrtc/bin",
		fmt.Sprintf("curl -sL %s -o ~/.owenrtc/bin/olcrtc", url),
		"chmod +x ~/.owenrtc/bin/olcrtc",
	}
	for _, s := range steps {
		if _, err := i.runner.Run(ctx, s); err != nil {
			return fmt.Errorf("install step %q: %w", s, err)
		}
	}
	return nil
}
