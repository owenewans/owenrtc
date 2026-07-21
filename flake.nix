{
  description = "owenrtc development environment - web panel for olcrtc";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { nixpkgs, ... }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      devShells = forAllSystems (system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go
              golangci-lint
              mage
              acme-sh
              wails
              pkg-config
              webkitgtk_4_1
              gtk3
            ];

            shellHook = ''
              export GOCACHE="''${XDG_CACHE_HOME:-$HOME/.cache}/go-build"
              export GOMODCACHE="''${XDG_CACHE_HOME:-$HOME/.cache}/go-mod"
            '';
          };
        });
    };
}
