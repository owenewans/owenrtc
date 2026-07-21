<div align="center">

```
                          __     
 ___ _    _____ ___  ____/ /_____
/ _ \ |/|/ / -_) _ \/ __/ __/ __/
\___/__,__/\__/_//_/_/  \__/\__/ 
```

minimal web panel for [olcrtc](https://github.com/openlibrecommunity/olcrtc)
<br>

</div>

stack       go 1.26+ / html / css / vanilla js / olcrtc / mage / gossh / wails
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
# install golang, mage, wails 
mage wails # or run
```

nix users:
```sh
git clone https://src.owenewans.org/owenrtc --recurse-submodules
cd owenrtc
nix develop
mage wails # or run
```

commands:
mage build  build panel binary
mage run    build + run panel on 127.0.0.1:8090
mage wails  build + run panel on wails
mage lint   golangci-lint
mage test   run tests

links:
docs        [olcrtc/docs](olcrtc/docs)
upstream    [openlibrecommunity/olcrtc](https://github.com/openlibrecommunity/olcrtc)

---
author      [owenewans.org](https://owenewans.org)
email       [owenewans@owenewans.org](mailto:owenewans@owenewans.org)
tg          [t.me/owenewans](https://t.me/owenewans)
