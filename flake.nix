{
  description = "A Nix-flake-based Go 1.27.3 development environment";

  inputs = {
    nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.2511";
    # Published compiler and matching builder; other tools retain their pins.
    nixpkgs-go-stable = {
      url = "https://api.flakehub.com/f/pinned/NixOS/nixpkgs/0.2605.1016949%2Brev-1f607aa79b6356ed45d3cd3c8079c196e79efa2c/01a17424-ae2d-7b20-83c4-cf14527294c3/source.tar.gz";
      flake = false;
    };
    nixpkgs-unstable.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1";
    go-tooling = {
      url = "github:NixOS/nixpkgs/8cec5ef1131078ea76689e84c362b0ed6a1fb6a1";
      flake = false;
    };
  };

  outputs = { nixpkgs, nixpkgs-go-stable, go-tooling, ... }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      forEachSupportedSystem = f:
        nixpkgs.lib.genAttrs supportedSystems (system:
          let
            pkgs = import nixpkgs {
              inherit system;
            };
          in
          f { inherit pkgs system; }
        );
    in
    {
      devShells = forEachSupportedSystem ({ pkgs, system }:
        let
          goPkgs = import nixpkgs-go-stable { inherit system; };
          go = goPkgs.go_1_27;
          buildGoModule = goPkgs.buildGo127Module;
          # Retain golangci-lint 2.15.0 with the matching upstream compiler.
          golangci-lint = goPkgs.callPackage "${go-tooling}/pkgs/by-name/go/golangci-lint/package.nix" {
            buildGo127Module = goPkgs.buildGo127Module;
          };
          mockgen = pkgs.mockgen.override { inherit buildGoModule; };
          stablePackages = with pkgs; [
            ginkgo
            gomarkdoc
            goperf
            gotools
            jq
            just
            yq-go
          ];
          goPackages = [ go mockgen golangci-lint ];
        in
        {
          default = pkgs.mkShell {
            packages = stablePackages ++ goPackages;

            env = {
              GOROOT = "${go}/share/go";
              GOTOOLCHAIN = "local";
            };
          };
        }
      );
    };
}
