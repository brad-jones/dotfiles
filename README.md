# dotfiles - managed by chezmoi

<https://github.com/twpayne/chezmoi>

These dotfiles have been developed on a Fedora system, other redhat based
distros may work but using anything else is certainly unsupported.

__UPDATE:__ Actually these days I do most of my work on company mandated Windows
laptop with WSL so thats what this is designed to setup but I still run a native
Fedora PC at home so this still works for that.

_NOTE: MacOS is not supported but could easily be if required in the future._

## Features

- `git` & `gpg` configuration for a Personal & Professional Identity.
- <https://github.com/gopasspw/gopass> for secret management.
- <https://github.com/99designs/aws-vault> _(with secret MFA sauce - connected to gopass)_ for managing AWS sessions.
- WSL support, a Fedora guest that automatically sets it's self up using this same repo - inception :wink:
- An ok PowerShell Core profile & Windows Terminal for when the job absolutely has to be done on Windows :joy:
- All the version managers you could want: `nodenv`, `goenv`, `rbenv`, `pyenv`, `tfenv`, `pkenv`, etc...

## Bootstrap

Just download `./setup/bin/setup_{{ YOUR OS }}_amd64` and run it.

## Apply Updates

Just execute `chezmoi update`.

## How this works

- We use `chezmoi` for anything we can do declaratively, static config files,
  what it was originally designed for.

- But for things we can not do 100% declaratively, eg: installing software.
  We have configured `chezmoi` through the `run_setup.bat.tmpl` &
  `run_setup.sh.tmpl` files to execute our golang `setup` tool.

- So when you execute `chezmoi update`, chezmoi will apply all the declarative
  updates to config files & the like first.

- And then it will always execute `~/.local/share/chezmoi/setup/bin/setup_{{ .chezmoi.os }}_amd64 chezmoi-apply`.
  _This is why we commit the compiled binaries to this repo._

- Our `setup` tool is designed to run in an [idempotent](https://en.wikipedia.org/wiki/Idempotence) manner.

- Executing the `setup` tool without any arguments will run the **Bootstrap** tasks.
  Which is also _idempotent_ so is safe to do run again if you want to totally start from scratch.