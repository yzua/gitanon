{
  description = "gitanon — anonymous git identity manager";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
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

            nativeCheckInputs = [ pkgs.git ];
            doCheck = true;
          };

          default = self.packages.${system}.gitanon;
        }
      );

      apps = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          gitanon = {
            type = "app";
            program = "${self.packages.${system}.gitanon}/bin/gitanon";
          };
          default = self.apps.${system}.gitanon;
        }
      );

      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              golangci-lint
              just
              git
            ];
          };
        }
      );
    };
}
