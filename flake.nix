{
  description = "A Nix-flake-based Go 1.27.3 development environment";

  inputs = {
    nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.2511";
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
      devShells = forEachSupportedSystem ({ pkgs, ... }:
        let
          # Import only Go recipes so the existing package set still supports Intel Macs.
          go = lintGo;
          lintGo = pkgs.callPackage "${go-tooling}/pkgs/development/compilers/go/1.27.nix" {
            buildGo127Module = lintBuildGoModule;
          };
          lintBuildGoModule = pkgs.buildGoModule.override { go = lintGo; };
          golangci-lint = pkgs.callPackage "${go-tooling}/pkgs/by-name/go/golangci-lint/package.nix" {
            buildGo127Module = lintBuildGoModule;
          };
          mockgen = pkgs.mockgen.override { buildGoModule = lintBuildGoModule; };
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
