// Package acme provisions SSL certificates for public IPs via acme.sh.
//
// Let's Encrypt supports IP certs with ~6 day validity, requiring
// full automation via acme.sh. Server mode uses this to serve the
// panel on :9443 with HTTPS.
package acme

import (
	"context"
	"fmt"

	"src.owenewans.org/owenrtc/internal/exec"
)

// Provisioner obtains SSL certs for a public IP address.
type Provisioner struct {
	runner exec.Runner
}

// New creates a provisioner with the given command runner.
func New(r exec.Runner) *Provisioner {
	return &Provisioner{runner: r}
}

// Provision obtains an SSL certificate for the given IP via acme.sh.
func (p *Provisioner) Provision(ctx context.Context, ip string) error {
	if ip == "" {
		return fmt.Errorf("empty ip")
	}
	steps := []string{
		"command -v acme.sh || curl https://get.acme.sh | sh",
		fmt.Sprintf("acme.sh --issue -d %s --standalone", ip),
	}
	for _, s := range steps {
		if _, err := p.runner.Run(ctx, s); err != nil {
			return fmt.Errorf("acme step %q: %w", s, err)
		}
	}
	return nil
}
