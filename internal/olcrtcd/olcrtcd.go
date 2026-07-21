// Package olcrtcd manages the olcrtc binary: download, run, stop, config.
package olcrtcd

import (
	"context"
	"fmt"
	"os/exec"
)

// Daemon wraps olcrtc binary management.
type Daemon struct {
	binary string
}

// New creates a daemon with the given binary path.
func New(binary string) *Daemon {
	if binary == "" {
		binary = "olcrtc"
	}
	return &Daemon{binary: binary}
}

// Start runs olcrtc with the given config path.
func (d *Daemon) Start(ctx context.Context, configPath string) (*exec.Cmd, error) {
	cmd := exec.CommandContext(ctx, d.binary, configPath)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start olcrtc: %w", err)
	}
	return cmd, nil
}

// Stop terminates a running olcrtc process.
func Stop(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	if err := cmd.Process.Kill(); err != nil {
		return fmt.Errorf("stop olcrtc: %w", err)
	}
	return nil
}
