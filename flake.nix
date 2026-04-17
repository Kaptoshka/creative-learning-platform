{
  description = "Web Platform for improving creative skills";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            # Go toolchain
            go_1_26
            gopls
            delve

            # Protobuf
            buf

            # Linting
            golangci-lint
            gotools # goimports

            yamllint
            statix
            deadnix
            nixpkgs-fmt
          ];

          shellHook = ''
            echo "Creative Learning Platform"
            echo "---------------------------"
            echo "- Go version: $(go version)"
          '';
        };
      }
    );
}
