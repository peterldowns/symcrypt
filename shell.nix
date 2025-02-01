# This is a shim that lets lorri work with flakes —  see
# https://github.com/target/lorri/issues/460#issuecomment-870901956
{ system ? builtins.currentSystem }:
(builtins.getFlake (toString ./.)).devShells.${system}.default
