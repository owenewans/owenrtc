//go:build !wails

// Package main is the browser entrypoint (default build).
// Runs the panel HTTP server on 127.0.0.1:8090 (or :9443 in server mode).
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"src.owenewans.org/owenrtc/internal/mode"
	"src.owenewans.org/owenrtc/internal/panel"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	m := mode.Detect()
	log.Printf("owenrtc: %s mode on :%d", m.Kind, m.Port)

	srv := panel.New(m)
	if err := srv.Run(ctx); err != nil {
		log.Fatalf("panel: %v", err)
	}
}
