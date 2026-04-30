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

        goCheck = pkgs.writeShellScriptBin "go-check" ''
          set -e

          golangci-lint run --config=${./.golangci.yml}
        '';

        nixCheck = pkgs.writeShellScriptBin "nix-check" ''
          set -e

          echo "==> Nix format (alejandra)"
          alejandra --check .

          echo "==> Static analysis"
          statix check .

          echo "==> Dead code"
          deadnix --fail .
        '';

        yamlCheck = pkgs.writeShellScriptBin "yaml-check" ''
          set -e
          yamllint -c ${./.yamllint.yml} .
        '';

        protoCheck = pkgs.writeShellScriptBin "proto-check" ''
          set -e
          buf lint ${./libs/protos}
        '';

        dockerCheck = pkgs.writeShellScriptBin "docker-check" ''
          set -e
          find . -type f \( -name "Dockerfile" -o -name "Dockerfile.dev" \) \
            -exec ${pkgs.hadolint}/bin/hadolint --config ${./.hadolint.yml} {} \;
        '';

        composeCheck = pkgs.writeShellScriptBin "compose-check" ''
          set -e
          ${pkgs.docker}/bin/docker compose -f docker-compose.yml config
          ${pkgs.docker}/bin/docker compose -f docker-compose.dev.yml config
        '';

      in
      {
        checks.default = pkgs.runCommand "ci-check" {} ''
          set -e

          echo "Running full CI checks..."

          ${goCheck}/bin/go-check
          ${nixCheck}/bin/nix-check
          ${yamlCheck}/bin/yaml-check
          ${protoCheck}/bin/proto-check
          ${dockerCheck}/bin/docker-check
          ${composeCheck}/bin/compose-check

          echo "All checks passed"
          touch $out
        '';

        checks.go = goCheck;
        checks.nix = nixCheck;
        checks.yaml = yamlCheck;
        checks.proto = protoCheck;
        checks.docker = dockerCheck;
        checks.compose = composeCheck;

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            # git toolchain
            pre-commit
            git

            # Go toolchain
            go_1_26
            gopls
            delve
            gofumpt
            goimports-reviser

            # Protobuf
            buf
            protoc-gen-go
            protoc-gen-go-grpc

            # Nix toolchain
            alejandra
            deadnix
            nixpkgs-fmt
            statix

            # Docker toolchain
            docker
            hadolint

            # Linting
            golangci-lint
            yamllint

            jq
            yq
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
