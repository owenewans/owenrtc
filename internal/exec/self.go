package exec

import (
	"context"
	"fmt"
	"os/exec"
)

// SelfRunner runs commands locally via os/exec.
type SelfRunner struct{}

// Run executes cmd on the local machine via sh -c.
func (SelfRunner) Run(ctx context.Context, cmd string) ([]byte, error) {
	out, err := exec.CommandContext(ctx, "sh", "-c", cmd).Output()
	if err != nil {
		return out, fmt.Errorf("self run %q: %w", cmd, err)
	}
	return out, nil
}
