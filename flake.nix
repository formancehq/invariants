{
  description = "A Nix-flake-based Go 1.27.3 development environment";

  inputs = {
    nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.2605";
    nixpkgs-unstable.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1";
    go-tooling = {
      url = "github:NixOS/nixpkgs/8cec5ef1131078ea76689e84c362b0ed6a1fb6a1";
      flake = false;
    };
  };

  outputs = { nixpkgs, go-tooling, ... }:
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
          # Use the upstream compiler and matching module builder for cache reuse.
          go = pkgs.go_1_27;
          buildGoModule = pkgs.buildGo127Module;
          # Stable has golangci-lint 2.12.2; preserve 2.15.0 and its cached dependencies.
          # The tooling package set no longer supports Intel macOS.
          golangci-lint =
            if system == "x86_64-darwin" then
              pkgs.callPackage "${go-tooling}/pkgs/by-name/go/golangci-lint/package.nix" { }
            else
              (import go-tooling { inherit system; }).golangci-lint;
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
