{
  description = "gitanon — anonymous git identity manager";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = {
          gitanon = pkgs.buildGoModule {
            pname = "gitanon";
            version = self.shortRev or self.dirtyShortRev or "dev";
            src = ./.;
            vendorHash = "sha256-7K17JaXFsjf163g5PXCb5ng2gYdotnZ2IDKk8KFjNj0=";

            ldflags = [
              "-s"
              "-w"
              "-X github.com/yzua/gitanon/cmd.Version=${self.shortRev or "dev"}"
            ];

            # Run tests during build
            doCheck = true;
          };

          default = self.packages.${system}.gitanon;
        };

        apps = {
          gitanon = {
            type = "app";
            program = "${self.packages.${system}.gitanon}/bin/gitanon";
          };
          default = self.apps.${system}.gitanon;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            golangci-lint
            just
            git
          ];
        };
      }
    );
}
