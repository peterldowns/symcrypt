# This Justfile contains rules/targets/scripts/commands that are used when
# developing. Unlike a Makefile, running `just <cmd>` will always invoke
# that command. For more information, see https://github.com/casey/just
#
#
# this setting will allow passing arguments through to tasks, see the docs here
# https://just.systems/man/en/chapter_24.html#positional-arguments
set positional-arguments

# print all available commands by default
help:
  just --list

# test go
test *args='./...':
  go test -race "$@"

# lint go + nix
lint *args:
  just lint-go
  just lint-nix

# lint go
lint-go:
  golangci-lint run --fix --config .golangci.yaml "$@"

# lint nix
lint-nix:
  find . -name '*.nix' | xargs nixpkgs-fmt

# go mod tidy
tidy:
  #!/usr/bin/env bash
  go mod tidy --go=1.20.0 --compat=1.20.0
