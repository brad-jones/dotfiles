#!/usr/bin/env bash
set -euxo pipefail;

sudo dnf update -y;
sudo dnf groupinstall -y "Development Tools";
sudo dnf install -y expect rng-tools wget jq tree bash-completion mlocate;
sudo dnf copr enable -y daftaupe/gopass;
sudo dnf install -y gopass;
chezmoiV="$(wget https://github.com/twpayne/chezmoi/releases/latest -O /dev/null 2>&1 | grep Location: | sed -r 's~^.*tag/v(.*?) \[.*~\1~g')";
sudo dnf install -y https://github.com/twpayne/chezmoi/releases/download/v$chezmoiV/chezmoi-$chezmoiV-x86_64.rpm;

rm -rf /tmp/vault-key;
rm -rf ~/.password-store;
rm -f /tmp/brad@bjc.id.au;
rm -f /tmp/brad.jones@xero.com;
rm -rf ~/.local/share/chezmoi;

git clone https://gitlab.com/brad-jones/vault-key.git /tmp/vault-key;
gpg --import /tmp/vault-key/private.pem;
expect -c "spawn gpg --edit-key \"Brad Jones (vault) <brad@bjc.id.au>\" trust quit; send \"5\ry\r\"; expect eof";
rm -rf /tmp/vault-key;

git clone https://github.com/brad-jones/vault.git ~/.password-store;
git --git-dir ~/.password-store/.git remote set-url origin git@github.com:brad-jones/vault.git;
gopass bin cp keys/ssh/brad@bjc.id.au ~/.ssh/brad@bjc.id.au;
gopass bin cp keys/ssh/brad.jones@xero.com ~/.ssh/brad.jones@xero.com

gopass bin cp keys/gpg/brad@bjc.id.au /tmp/brad@bjc.id.au;
gpg --import /tmp/brad@bjc.id.au;
expect -c "spawn gpg --edit-key \"Brad Jones <brad@bjc.id.au>\" trust quit; send \"5\ry\r\"; expect eof";
rm -f /tmp/brad@bjc.id.au;

gopass bin cp keys/gpg/brad.jones@xero.com /tmp/brad.jones@xero.com;
gpg --import /tmp/brad.jones@xero.com;
expect -c "spawn gpg --edit-key \"Brad Jones <brad.jones@xero.com>\" trust quit; send \"5\ry\r\"; expect eof";
rm -f /tmp/brad.jones@xero.com;

chezmoi init https://github.com/brad-jones/dotfiles.git;
git --git-dir "$(chezmoi source-path)/.git" remote set-url origin git@github.com:brad-jones/dotfiles.git;
chezmoi apply --debug;

sudo shutdown -r now;
