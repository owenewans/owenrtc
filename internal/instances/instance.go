// Package instances manages olcrtc server/client instances with limits.
//
// Limits (traffic, speed) are enforced at the egress SOCKS proxy,
// not at the olcrtc core. If an outbound SOCKS is specified, all
// traffic goes through it before hitting the internet.
package instances

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Provider is the auth provider.
type Provider string

const (
	Jitsi    Provider = "jitsi"
	Telemost Provider = "telemost"
	WBStream Provider = "wbstream"
)

// Transport is the data transport.
type Transport string

const (
	DataChannel  Transport = "datachannel"
	VP8Channel   Transport = "vp8channel"
	SEIChannel   Transport = "seichannel"
	VideoChannel Transport = "videochannel"
)

// Limits are enforced at the egress socks, not at olcrtc core.
type Limits struct {
	TrafficLimit int64 `json:"traffic_limit"` // bytes, 0 = unlimited
	SpeedLimit   int64 `json:"speed_limit"`   // bytes/sec, 0 = unlimited
}

// OutboundSOCKS is an optional upstream SOCKS5 proxy for egress traffic.
type OutboundSOCKS struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user,omitempty"`
	Pass string `json:"pass,omitempty"`
}

// Instance is a running olcrtc server or client.
type Instance struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Provider  Provider       `json:"provider"`
	Transport Transport      `json:"transport"`
	RoomID    string         `json:"room_id"`
	Key       string         `json:"key"`
	Limits    Limits         `json:"limits"`
	Outbound  *OutboundSOCKS `json:"outbound,omitempty"`
	CreatedAt int64          `json:"created_at"` // unix timestamp
	Status    string         `json:"status"`
}

// ServerYAML is the server-side olcrtc config written to disk.
type ServerYAML struct {
	Mode   string `yaml:"mode"`
	Auth   struct {
		Provider string `yaml:"provider"`
	} `yaml:"auth"`
	Room struct {
		ID string `yaml:"id"`
	} `yaml:"room"`
	Crypto struct {
		Key string `yaml:"key"`
	} `yaml:"crypto"`
	Net struct {
		Transport string `yaml:"transport"`
		DNS       string `yaml:"dns"`
	} `yaml:"net"`
	Data string `yaml:"data"`
}

// ClientYAML is the client-side olcrtc config.
type ClientYAML struct {
	Mode   string `yaml:"mode"`
	Auth   struct {
		Provider string `yaml:"provider"`
	} `yaml:"auth"`
	Room struct {
		ID string `yaml:"id"`
	} `yaml:"room"`
	Crypto struct {
		Key string `yaml:"key"`
	} `yaml:"crypto"`
	Net struct {
		Transport string `yaml:"transport"`
		DNS       string `yaml:"dns"`
	} `yaml:"net"`
	Socks struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"socks"`
	Data string `yaml:"data"`
}

// NewID generates a random hex ID.
func NewID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// NewKey generates a 32-byte hex crypto key.
func NewKey() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Now returns current unix timestamp.
func Now() int64 { return time.Now().Unix() }

func storePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	dir := filepath.Join(home, ".owenrtc")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}
	return filepath.Join(dir, "instances.json"), nil
}

// LoadAll reads all instances from disk.
func LoadAll() ([]Instance, error) {
	p, err := storePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return []Instance{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read instances: %w", err)
	}
	var list []Instance
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("parse instances: %w", err)
	}
	return list, nil
}

// SaveAll writes all instances to disk.
func SaveAll(list []Instance) error {
	p, err := storePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal instances: %w", err)
	}
	if err := os.WriteFile(p, data, 0o600); err != nil {
		return fmt.Errorf("write instances: %w", err)
	}
	return nil
}

// Add persists a new instance with generated ID and timestamp.
func Add(inst *Instance) error {
	if inst.ID == "" {
		inst.ID = NewID()
	}
	if inst.CreatedAt == 0 {
		inst.CreatedAt = Now()
	}
	if inst.Status == "" {
		inst.Status = "stopped"
	}
	list, err := LoadAll()
	if err != nil {
		return err
	}
	list = append(list, *inst)
	return SaveAll(list)
}

// Remove deletes an instance by ID.
func Remove(id string) error {
	list, err := LoadAll()
	if err != nil {
		return err
	}
	out := list[:0]
	for _, i := range list {
		if i.ID != id {
			out = append(out, i)
		}
	}
	return SaveAll(out)
}

// Find returns an instance by ID.
func Find(id string) (*Instance, error) {
	list, err := LoadAll()
	if err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].ID == id {
			return &list[i], nil
		}
	}
	return nil, fmt.Errorf("instance %s not found", id)
}
