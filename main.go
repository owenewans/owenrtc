//go:build wails

// Package main is the wails desktop entrypoint.
// Embeds web/ as assets and binds the API service to frontend.
// Frontend calls Go methods via wailsjs bindings, not HTTP.
package main

import (
	"context"
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"src.owenewans.org/owenrtc/internal/panel"
)

//go:embed all:web
var assets embed.FS

func main() {
	api := panel.NewAPI()

	err := wails.Run(&options.App{
		Title:  "owenrtc",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: api.Startup,
		Bind: []any{
			api,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = context.Background
}
