{
  description = "A basic gomod2nix flake";
  inputs =  {
  ewxfront.url = "github:kyp0717/ewxfront?ref=test-react";
  nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  flake-utils.url = "github:numtide/flake-utils";
  gomod2nix.url = "github:nix-community/gomod2nix";
  gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  gomod2nix.inputs.flake-utils.follows = "flake-utils";
  };

  outputs =
    { self,
      nixpkgs,
      flake-utils,
      gomod2nix,
      ewxfront,
    }:
    (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.callPackage ./default.nix {
          inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
        };
        devShells.default = pkgs.callPackage ./shell.nix {
          inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
        };

        devShells.nvim = ewxfront.devShells.${system}.default;

      }
    ));
}
