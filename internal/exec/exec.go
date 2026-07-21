// Package exec runs commands either locally (self) or over SSH.
//
// Each command is a standalone string with an execution method.
// Method "self" runs locally, "ssh" runs over an SSH connection.
package exec

import (
	"context"
	"fmt"
)

// Method is how a command is executed.
type Method string

const (
	Self Method = "self" // run locally
	SSH  Method = "ssh"  // run over ssh
)

// Command is a standalone command with its execution method.
type Command struct {
	Cmd    string
	Method Method
}

// Runner executes commands and returns combined output.
type Runner interface {
	Run(ctx context.Context, cmd string) ([]byte, error)
}

// ErrNoRunner is returned when no runner is available for the given method.
var ErrNoRunner = fmt.Errorf("no runner for method")

// Pick returns the right runner for the given method.
func Pick(m Method, self, ssh Runner) (Runner, error) {
	switch m {
	case Self:
		return self, nil
	case SSH:
		return ssh, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrNoRunner, m)
	}
}
