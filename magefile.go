//go:build mage

// Package main is the magefile for owenrtc.
//
// build   compile the panel binary
// run     build + run panel on 127.0.0.1:8090
// wails   build + run as wails desktop app
// check   build + vet + lint + test (pre-commit)
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	binary   = "owenrtc"
	buildDir = "build"
)

func init() {
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
}

// Build compiles the owenrtc panel binary.
func Build() error {
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return fmt.Errorf("mkdir build: %w", err)
	}
	if err := sh.RunV("go", "build", "-o", filepath.Join(buildDir, binary), "./cmd/owenrtc"); err != nil {
		return fmt.Errorf("build: %w", err)
	}
	return nil
}

// Run builds and runs the panel on 127.0.0.1:8090.
func Run() error {
	mg.Deps(Build)
	return sh.RunV(filepath.Join(buildDir, binary))
}

// Wails builds and runs the panel as a wails desktop app.
//
// On Wayland, webkit2gtk's DMABUF renderer breaks webview input and
// interactivity (logo renders but clicks/typing dead). We force the x11
// backend and disable DMABUF, and disable compositing as extra safety.
// See wailsapp/wails#3893.
func Wails() error {
	if _, err := exec.LookPath("wails"); err != nil {
		fmt.Println("wails not found - install options:")
		fmt.Println("  go:    go install github.com/wailsapp/wails/v2/cmd/wails@latest")
		fmt.Println("  arch:  pacman -S webkit2gtk gtk3")
		fmt.Println("  apt:   see https://wails.io/docs/guides/linux")
		fmt.Println("  then:  wails dev -tags wails")
		fmt.Println("")
		fmt.Println("or use browser mode: mage run -> http://127.0.0.1:8090")
		return nil
	}
	env := map[string]string{
		"GDK_BACKEND":                  "x11",
		"WEBKIT_DISABLE_DMABUF_RENDERER": "1",
		"WEBKIT_DISABLE_COMPOSITING_MODE": "1",
	}
	return sh.RunWithV(env, "wails", "dev", "-tags", "wails webkit2_41")
}

// Lint runs golangci-lint.
func Lint() error {
	return sh.RunV("golangci-lint", "run")
}

// Vet runs go vet.
func Vet() error {
	return sh.RunV("go", "vet", "./...")
}

// Test runs all tests with race detector.
func Test() error {
	return sh.RunV("go", "test", "-race", "-count=1", "./...")
}

// Check runs build + vet + lint + test (pre-commit).
func Check() error {
	mg.Deps(Build, Vet, Lint, Test)
	return nil
}

// Tidy tidies and verifies go modules.
func Tidy() error {
	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}
	return sh.RunV("go", "mod", "verify")
}

// Clean removes build artifacts.
func Clean() error {
	return os.RemoveAll(buildDir)
}
