# Brads DotFiles

This repo, represents Brads home directory dotfiles for both _*nix_ based
systems & Windows ones. It is a Go project which gets compiled into a single
binary, download & run the binary to install Brad's dotfiles. _Brad does not
use MacOS... yet_

> Fedora is my primary distro, other Redhat based distros may work but YMMV.

## Quick Start

### *nix

```
curl -fsSL https://github.com/brad-jones/dotfiles/releases/latest/download/install.sh | sh
```

### Windows

```
iwr https://github.com/brad-jones/dotfiles/releases/latest/download/install.ps1 -useb | iex
```

These one liners are designed to get you up and running on a fresh system,
quickly & easily. The scripts are designed to be small, simple and auditable.
_You should never just execute any arbitrary code from the internet!_

These scripts have the following config variables, which can be overridden via
corresponding environment variables:

- `DOTFILES_VERSION`: The default value is etched into the script by the pipeline at release time.
- `DOTFILES_INSTALL_DIR`: Defaults to `~/.local/bin`.
- `DOTFILES_DOWNLOAD_URL`: Constructed from `DOTFILES_VERSION`.
- `DOTFILES_OUTPUT`: Defaults to `$DOTFILES_INSTALL_DIR/dotfiles[.exe]`
- `DOTFILES_SIGNERS`: Defaults to `0x33dc7b56c2be6175e1ad17e31f003f55943fa4ce`
- `DOTFILES_RUN_AFTER_INSTALL`: Default to `true`

### CodeNotary

All artifacts produced by the pipeline are signed with <https://codenotary.io> -
A blockchain based digital code signing service.

The above one liners will download the executable and then verify it against
the CodeNotary API before executing it for you.

To do this manually just download the correct executable from <https://github.com/brad-jones/dotfiles/releases>.
Then open <https://authenticate.codenotary.io> in a browser & drop the downloaded file into the page.
This will hash the file locally & send the hash _(not the file)_ off to the blockchain for verification.

My signing ID is:

```
0x33dc7b56c2be6175e1ad17e31f003f55943fa4ce
```

## Features

- Zero dependency installation & 100% automated setup.

- 100% idempotent, ie: safe to run many times over and get the same results.
  <https://en.wikipedia.org/wiki/Idempotence>

- Automatically updates it's self on login / this repo automatically updates
  through [dependabot](https://dependabot.com/) _like_ workflows.

- On Windows based systems it will create a WSL instance and then automatically
  run the Linux version of the binary (which is embedded into the Windows one)
  and apply the setup inside the VM as well - inception :wink:

- Easily apply rollbacks by simply downloading & running a previously released
  version of the self-contained dotfiles binary.

- `ssh` / `gpg` / `git` configuration for a Personal & Professional Identity.

- Uses [gopass](https://www.gopass.pw/) for secret management.

  - I am thinking about replacing this with something else.
    I find the GPG/team management features are overkill.
    These are my personal secrets that I will never share with anyone else.
    A symmetric encryption based password manager would be preferred.

- [aws-vault](https://github.com/99designs/aws-vault) _(with secret MFA sauce -
 connected to gopass)_ for managing AWS sessions.

- On Linux we install all the version managers you could want:
  `nodenv`, `goenv`, `rbenv`, `pyenv`, `tfenv`, `pkenv`, etc...

- On Windows we just use [scoop](https://scoop.sh/) and install the latest
  version of such tools.

- An ok PowerShell profile & Windows Terminal for when the job absolutely has
  to be done on Windows :joy:

- And much more...

## Why this custom thing & not some off the shelf dotfile manager?

There was a time where I only cared for Linux, I started out with a
[chezmoi](https://www.chezmoi.io/) repo for my Linux dotfiles and this
worked very well.

But then I got a new job at [Xero](https://www.xero.com/) and had to make do
with a Windows laptop.

While [chezmoi](https://www.chezmoi.io/) does work on Windows, it's got some
rough edges & to create a single [chezmoi](https://www.chezmoi.io/) repo that
works well for both Linux & Windows became a challenge.

Maybe thats the answer, to create two separate repositories but then you would
have some duplicated stuff & I didn't want that either.

The other main concern was that [chezmoi](https://www.chezmoi.io/) focuses on a
declarative approach and discourages any imperative setup like the installing of
software.

While I see why they have made this choice & on *nix based systems you can
easily work with-in the declarative guard rails, on Windows I found this more
challenging due to things like the Registry urgh...

And I couldn't help but feel that it's a bit of a cop out to just discard that
part of ones environment setup. For me I wanted something 100% automated.
_(Or as close as possible to this that I can get.)_
