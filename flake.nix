{
  description = "A Nix-flake-based Go 1.26 development environment";

  inputs = {
    nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.2511";
    nixpkgs-unstable.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1";
  };

  outputs = { nixpkgs, nixpkgs-unstable, ... }:
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
            pkgs-unstable = import nixpkgs-unstable {
              inherit system;
            };
          in
          f { inherit pkgs pkgs-unstable system; }
        );
    in
    {
      devShells = forEachSupportedSystem ({ pkgs, pkgs-unstable, ... }:
        let
          stablePackages = with pkgs; [
            ginkgo
            go_1_26
            gomarkdoc
            goperf
            gotools
            jq
            just
            mockgen
            yq-go
          ];
          unstablePackages = with pkgs-unstable; [
            golangci-lint
          ];
        in
        {
          default = pkgs.mkShell {
            packages = stablePackages ++ unstablePackages;

            env = {
              GOROOT = "${pkgs.go_1_26}/share/go";
            };
          };
        }
      );
    };
}
