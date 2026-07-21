<div align="center">

```text
                          __     
 ___ _    _____ ___  ____/ /_____
/ _ \ |/|/ / -_) _ \/ __/ __/ __/
\___/__,__/\__/_//_/_/  \__/\__/ 
```

minimal web panel for [olcrtc](https://github.com/openlibrecommunity/olcrtc)
<br>
</div>

stack       go 1.26+ / html / css / vanilla js / olcrtc / mage / gossh
license     wtfpl

features:
- start/stop srv and cnc from browser
- edit yaml configs with presets
- generate and store crypto keys
- build olcrtc:// uris and sub.md subscriptions
- logs
- limitations

quick start:
```sh
git clone https://src.owenewans.org/owenrtc --recurse-submodules
cd owenrtc
mage run
```

nix users:
```sh
git clone https://src.owenewans.org/owenrtc --recurse-submodules
cd owenrtc
nix develop
mage run
```

wails (optional, desktop app):
wails v2 needs webkit2gtk-4.0, not in nixpkgs unstable.
install manually: go install github.com/wailsapp/wails/v2/cmd/wails@latest
then: wails dev -tags wails

commands:
mage build  build panel binary
mage run    build + run panel on 127.0.0.1:8090
mage wails  build + run as wails desktop (needs manual wails install)
mage lint   golangci-lint
mage test   run tests

links:
docs        [olcrtc/docs](olcrtc/docs)
upstream    [openlibrecommunity/olcrtc](https://github.com/openlibrecommunity/olcrtc)

---
author      [owenewans.org](https://owenewans.org)
email       [owenewans@owenewans.org](mailto:owenewans@owenewans.org)
tg          [t.me/owenewans](https://t.me/owenewans)
