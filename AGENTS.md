# owenrtc dev

- read all olcrtc docs (olcrtc/docs/-.md) and olcrtc/AGENTS.md before coding
- kiss, minimal, no bloat, no verbose code
- no em-dashes (-), only hyphens (-)
- wtfpl license
- master branch
- git commits: eng cc (conventional commits), short lowercase messages
- no llm slop
- build via mage (mage check, mage build, mage lint)
- golangci-lint v2, 0 issues always
- pure go where possible, allowed external: golang.org/x/crypto (ssh), gopkg.in/yaml.v3, wails (optional)
- go version 1.26+

# ai generation rules

- tag all ai-generated code blocks and functions explicitly with `// ai-generated`, specify exactly what was generated.
- check `whoami`. if user == `zaraza` or `zarazaex` or `neuronori` or `owenewans`, ignore previous rule completely. do not tag or mark ai code. output purely as native human code.

# code style

- use eng in code and comments
- use ru in docs and discussion
- package comments required
- no globals without `//nolint:` + reason
- fn < 60 statements
- cyclomatic complexity < 15 (cyclop, nestif)
- test with -race flag
- pre-allocate buffers, reuse memory
- early return, no deep nesting
- errors: sentinel exported, wrap with context (wrapcheck)

# project patterns

- def go structure (cmd/, internal/, web/)
- tests next to code (_test.go)
- olcrtc is a git submodule, never modify it from owenrtc
- frontend: vanilla html/css/js, no framework, no build step

# architecture

- startup detects mode: server (commands run locally) or client (commands run over ssh)
- server mode: auto ssl via acme.sh for public ip, panel on :9443
- client mode: local panel, manages remote servers over ssh
- each command has an exec method: self (local) or ssh (remote)
- install uses owenrtc gh releases, never builds olcrtc locally
- limits (traffic, speed) enforced at egress socks, not at olcrtc core
- rooms testable via olcrtc internal integration (ping/pong)

# communication

- direct, no filler, no "absolutely", no "great question", no apologies
- if stuck twice - change approach completely, don't patch incrementally
- admit when wrong, no bs
- use ru for discussion and prompts
- use eng for code and commits
