#!/usr/bin/env bash
set -euxo pipefail;

# Install some tools
sudo dnf update -y;
sudo dnf groupinstall -y "Development Tools";
sudo dnf install -y expect rng-tools wget jq tree bash-completion mlocate;
sudo dnf copr enable -y daftaupe/gopass;
sudo dnf install -y gopass;

# Install chezmoi
chezmoiV="$(wget https://github.com/twpayne/chezmoi/releases/latest -O /dev/null 2>&1 | grep Location: | sed -r 's~^.*tag/v(.*?) \[.*~\1~g')";
sudo dnf install -y https://github.com/twpayne/chezmoi/releases/download/v$chezmoiV/chezmoi-$chezmoiV-x86_64.rpm;

# Ensure this script is Idempotent
rm -rf /tmp/vault-key;
rm -rf ~/.password-store;
rm -f /tmp/brad@bjc.id.au;
rm -f ~/.ssh/brad@bjc.id.au;
rm -rf ~/.local/share/chezmoi;
rm -f /tmp/brad.jones@xero.com;
rm -f ~/.ssh/brad@bjc.id.au.pub;
rm -f ~/.ssh/brad.jones@xero.com;
rm -f ~/.ssh/brad.jones@xero.com.pub;

# Install the GPG key from gitlab that is used to decrypt my gopass vault
git config --global credential.helper 'cache';
git clone https://gitlab.com/brad-jones/vault-key.git /tmp/vault-key;
gpg --import /tmp/vault-key/private.pem;
expect -c "spawn gpg --edit-key \"Brad Jones (vault) <brad@bjc.id.au>\" trust quit; send \"5\ry\r\"; expect eof";
rm -rf /tmp/vault-key;

# Install the gopass vault from github. To unlock the vault we need to know
# 3 things (gitlab password, github password & the key passphrase).
git clone https://github.com/brad-jones/vault.git ~/.password-store;
git --git-dir ~/.password-store/.git remote set-url origin git@github.com:brad-jones/vault.git;

# Install my personal and work SSH keys
gopass bin cp keys/ssh/brad@bjc.id.au ~/.ssh/brad@bjc.id.au;
gopass bin cp keys/ssh/brad.jones@xero.com ~/.ssh/brad.jones@xero.com

# Install my personal GPG key
gopass bin cp keys/gpg/brad@bjc.id.au /tmp/brad@bjc.id.au;
gpg --import /tmp/brad@bjc.id.au;
expect -c "spawn gpg --edit-key \"Brad Jones <brad@bjc.id.au>\" trust quit; send \"5\ry\r\"; expect eof";
rm -f /tmp/brad@bjc.id.au;

# Install my work GPG key
gopass bin cp keys/gpg/brad.jones@xero.com /tmp/brad.jones@xero.com;
gpg --import /tmp/brad.jones@xero.com;
expect -c "spawn gpg --edit-key \"Brad Jones <brad.jones@xero.com>\" trust quit; send \"5\ry\r\"; expect eof";
rm -f /tmp/brad.jones@xero.com;

# Install my dotfiles
chezmoi init https://github.com/brad-jones/dotfiles.git;
git --git-dir "$(chezmoi source-path)/.git" remote set-url origin git@github.com:brad-jones/dotfiles.git;
chezmoi apply --debug;

# Reboot to make sure things like kernels are updated etc
sudo shutdown -r now;
