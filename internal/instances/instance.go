// Package instances manages olcrtc server/client instances with limits.
//
// Limits (traffic, speed) are enforced at the egress SOCKS proxy,
// not at the olcrtc core. If an outbound SOCKS is specified, all
// traffic goes through it before hitting the internet.
package instances

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
	Key       string         `json:"-"`
	Limits    Limits         `json:"limits"`
	Outbound  *OutboundSOCKS `json:"outbound,omitempty"`
	CreatedAt int64          `json:"created_at"` // unix timestamp
	Status    string         `json:"status"`
}
