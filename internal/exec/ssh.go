package exec

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHConfig holds SSH connection details.
type SSHConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Key      []byte
}

// SSHRunner runs commands over SSH.
type SSHRunner struct {
	cfg SSHConfig
}

// NewSSHRunner creates a new SSH runner.
func NewSSHRunner(cfg SSHConfig) (*SSHRunner, error) {
	if cfg.Port == 0 {
		cfg.Port = 22
	}
	return &SSHRunner{cfg: cfg}, nil
}

// Run executes cmd over SSH.
func (r *SSHRunner) Run(ctx context.Context, cmd string) ([]byte, error) {
	auth, err := r.authMethod()
	if err != nil {
		return nil, fmt.Errorf("ssh auth: %w", err)
	}

	cfg := &ssh.ClientConfig{
		User:            r.cfg.User,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := net.JoinHostPort(r.cfg.Host, fmt.Sprintf("%d", r.cfg.Port))
	client, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return nil, fmt.Errorf("ssh dial %s: %w", addr, err)
	}
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("ssh session: %w", err)
	}
	defer sess.Close()

	type result struct {
		out []byte
		err error
	}
	ch := make(chan result, 1)
	go func() {
		var buf bytes.Buffer
		sess.Stdout = &buf
		sess.Stderr = &buf
		err := sess.Run(cmd)
		ch <- result{buf.Bytes(), err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-ch:
		if res.err != nil {
			return res.out, fmt.Errorf("ssh run %q: %w", cmd, res.err)
		}
		return res.out, nil
	}
}

func (r *SSHRunner) authMethod() (ssh.AuthMethod, error) {
	if len(r.cfg.Key) > 0 {
		signer, err := ssh.ParsePrivateKey(r.cfg.Key)
		if err != nil {
			return nil, fmt.Errorf("parse key: %w", err)
		}
		return ssh.PublicKeys(signer), nil
	}
	// many servers disable "password" auth and require keyboard-interactive.
	// provide both methods so ssh picks whichever the server allows.
	return ssh.PasswordCallback(func() (string, error) {
		return r.cfg.Password, nil
	}), nil
}
