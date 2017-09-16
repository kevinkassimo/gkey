# GOKEY password keeper

__GOKEY__ is a minimal password keeper that uses AES encryption to save your password.
It is a Go practice project. Code may need further improvements

[![asciicast](https://asciinema.org/a/jJn3iYozWRoMwJ7YTT3TLOZms.png)](https://asciinema.org/a/jJn3iYozWRoMwJ7YTT3TLOZms)

## Install
Run `install.sh`. Compiled `gokey` would be installed at `$GOPATH/bin`. Make sure that your `$GOPATH` is set.
User information would be stored in `~/.gokey_store`. You can copy this folder to any other machines so that the records could also be queried.

## Usage
Run `gokey -h` after installation to see the detailed usage.