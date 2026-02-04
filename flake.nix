{
  description = "Go DevShell";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs =
    { nixpkgs, ... }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs {
        inherit system;
      };
      llvm = pkgs.llvmPackages;
    in
    {
      devShells.${system}.default =
        pkgs.mkShell.override
          {
            stdenv = llvm.libcxxStdenv;
          }
          {
            packages = with pkgs; [
              # Utilities
              nixd
              nil
              package-version-server
              cloc

              # Go
              go
            ];

            shellHook = ''
              echo "======== Go DevShell ========"
            '';
          };

    };
}
