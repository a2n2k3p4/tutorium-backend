{
  pkgs,
  lib,
  config,
  inputs,
  ...
}: {
  # https://devenv.sh/basics/
  # env.GREET = "devenv";

  # https://devenv.sh/packages/
  # packages = [ pkgs.git ];
  packages = with pkgs; [codespell go-swag];

  # https://devenv.sh/languages/
  # languages.rust.enable = true;
  languages.go.enable = true;

  # https://devenv.sh/processes/
  # processes.cargo-watch.exec = "cargo-watch";

  # https://devenv.sh/services/
  # services.postgres.enable = true;
  services = {
    postgres = {
      enable = true;
      package = pkgs.postgresql_17;
      initialDatabases = [
        {
          name = config.env.DB_NAME;
          pass = config.env.DB_PASSWORD;
          user = config.env.DB_USER;
        }
      ];
      listen_addresses = config.env.DB_HOST;
      port = lib.toInt config.env.DB_PORT;
    };
    minio = {
      enable = true;
      accessKey = config.env.MINIO_ACCESS_KEY;
      buckets = [config.env.MINIO_BUCKET];
      listenAddress = config.env.MINIO_ENDPOINT;
      secretKey = config.env.MINIO_SECRET_KEY;
    };
  };

  # https://devenv.sh/scripts/
  # scripts.hello.exec = ''
  #   echo hello from $GREET
  # '';

  # enterShell = ''
  #   hello
  #   git --version
  # '';

  # https://devenv.sh/tasks/
  # tasks = {
  #   "myproj:setup".exec = "mytool build";
  #   "devenv:enterShell".after = [ "myproj:setup" ];
  # };

  # https://devenv.sh/tests/
  # enterTest = ''
  #   echo "Running tests"
  #   git --version | grep --color=auto "${pkgs.git.version}"
  # '';

  # https://devenv.sh/git-hooks/
  # git-hooks.hooks.shellcheck.enable = true;

  # See full reference at https://devenv.sh/reference/options/

  dotenv.enable = true;
}
