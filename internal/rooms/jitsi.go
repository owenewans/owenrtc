// Package rooms loads jitsi instances and tests room connectivity.
package rooms

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// JitsiInstances is the yaml structure of jitsi.instances.yaml.
type JitsiInstances struct {
	Instances []string `yaml:"instances"`
}

// LoadJitsiInstances reads jitsi.instances.yaml from the olcrtc submodule.
func LoadJitsiInstances(submoduleDir string) ([]string, error) {
	p := filepath.Join(submoduleDir, "docs", "examples", "jitsi.instances.yaml")
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("read jitsi instances: %w", err)
	}
	var ji JitsiInstances
	if err := yaml.Unmarshal(data, &ji); err != nil {
		return nil, fmt.Errorf("parse jitsi instances: %w", err)
	}
	return ji.Instances, nil
}

// TestResult is the outcome of a room connectivity test.
type TestResult struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// TestRoom joins a room and sends ping/pong via olcrtc internal integration.
func TestRoom(_ context.Context, _, _, _ string) (*TestResult, error) {
	// TODO: use olcrtc internal integration to join room and ping/pong
	return &TestResult{OK: false, Message: "not implemented"}, nil
}
