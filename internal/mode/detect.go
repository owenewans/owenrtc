// Package mode detects whether owenrtc runs on a server or client.
package mode

import (
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// Kind is the runtime mode.
type Kind string

const (
	Server Kind = "server" // commands run locally, ssl on :9443
	Client Kind = "client" // commands run over ssh, local panel
)

// Mode is the detected runtime configuration.
type Mode struct {
	Kind     Kind
	Port     int
	PublicIP string
}

// Detect determines the runtime mode by checking for a public IP.
func Detect() Mode {
	ip := publicIP()
	if ip != "" && isLocalIP(ip) {
		return Mode{Kind: Server, Port: 9443, PublicIP: ip}
	}
	return Mode{Kind: Client, Port: 8090}
}

func publicIP() string {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://icanhazip.com")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(body))
}

func isLocalIP(ip string) bool {
	ifaces, err := net.Interfaces()
	if err != nil {
		return false
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if strings.HasPrefix(addr.String(), ip+"/") {
				return true
			}
		}
	}
	return false
}
