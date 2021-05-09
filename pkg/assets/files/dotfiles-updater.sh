#!/usr/bin/env bash
set -euxo pipefail;

# Source our bashrc file
[ -f ~/.bashrc ] && set +u && . ~/.bashrc && set -u;

# Update the entire system
sudo dnf update -y;

# Install Docker (Podman)
# ------------------------------------------------------------------------------
sudo dnf install -y podman;
sudo dnf reinstall -y shadow-utils;
if ! [ -x "/usr/bin/docker" ]; then
	sudo ln -s "/usr/bin/podman" "/usr/bin/docker";
fi

# Install awscli
# ------------------------------------------------------------------------------
if ! [ -x "$HOME/.local/bin/aws" ]; then
    rm -rf ~/.local/aws-cli ~/.local/bin/aws;
fi
tmpFolderAWS="/tmp/$(uuidgen)";
mkdir -p $tmpFolderAWS;
function finish {
	rm -rf $tmpFolderAWS;
}
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "$tmpFolderAWS/awscliv2.zip";
unzip "$tmpFolderAWS/awscliv2.zip" -d "$tmpFolderAWS/extracted";
$tmpFolderAWS/extracted/aws/install -i ~/.local/aws-cli -b ~/.local/bin --update;
~/.local/bin/aws --version;

# Install Deno
# ------------------------------------------------------------------------------
if [ -z "$(command -v asdf)" ]; then
    rm -rf ~/.asdf;
	sudo dnf install -y curl git;
	git clone https://github.com/asdf-vm/asdf.git ~/.asdf;
	cd ~/.asdf && git checkout "$(git describe --abbrev=0 --tags)" && cd -;
    . "$HOME/.asdf/asdf.sh";
	asdf plugin-add deno https://github.com/asdf-community/asdf-deno.git;
fi
asdf update;
asdf plugin update deno;
asdf install deno latest;
asdf global deno $(asdf latest deno);

# Install Dartlang
# ------------------------------------------------------------------------------
if [ -z "$(command -v dart)" ]; then
    rm -rf ~/.dart && mkdir ~/.dart;
fi
tmpFolderDart="/tmp/$(uuidgen)";
mkdir -p $tmpFolderDart;
function finish {
	rm -rf $tmpFolderDart;
}
curl "https://storage.googleapis.com/dart-archive/channels/stable/release/latest/sdk/dartsdk-linux-x64-release.zip" -o "$tmpFolderDart/dartsdk.zip";
unzip "$tmpFolderDart/dartsdk.zip" -d "$tmpFolderDart/extracted";
dartV="$(cat "$tmpFolderDart/extracted/dart-sdk/version")";
if ! [ -d "$HOME/.dart/$dartV" ]; then
	mv "$tmpFolderDart/extracted/dart-sdk" "$HOME/.dart/$dartV";
	ln -s "$HOME/.dart/$dartV" "$HOME/.dart/current";
fi

# Install Golang
# ------------------------------------------------------------------------------
if [ -z "$(command -v goenv)" ]; then
    rm -rf ~/.goenv;
    git clone https://github.com/syndbg/goenv.git ~/.goenv;
    cd ~/.goenv && src/configure && make -C src && cd -;
fi
cd ~/.goenv && git pull && cd -;
goV="$(echo "$(~/.goenv/bin/goenv install --list)" | awk '{$1=$1};1' | grep '^[0-9]' | grep -v '[a-zA-Z]' | sed '/-/!{s/$/_/}' | sort -V | sed 's/_$//' | tail -n1)";
~/.goenv/bin/goenv install -s $goV;
~/.goenv/bin/goenv global $goV;
~/.goenv/bin/goenv rehash;

# Install Nodejs
# ------------------------------------------------------------------------------
if [ -z "$(command -v nodenv)" ]; then
    rm -rf ~/.nodenv;
    git clone https://github.com/nodenv/nodenv.git ~/.nodenv;
    mkdir -p ~/.nodenv/plugins;
    git clone https://github.com/nodenv/node-build.git ~/.nodenv/plugins/node-build;
    git clone https://github.com/nodenv/nodenv-update.git ~/.nodenv/plugins/nodenv-update;
    git clone https://github.com/nodenv/node-build-update-defs.git ~/.nodenv/plugins/node-build-update-defs;
    cd ~/.nodenv && src/configure && make -C src && cd -;
fi
~/.nodenv/bin/nodenv update;
nodeV="$(echo "$(~/.nodenv/bin/nodenv install --list)" | awk '{$1=$1};1' | grep '^[0-9]' | grep -v '[a-zA-Z]' | sed '/-/!{s/$/_/}' | sort -V | sed 's/_$//' | tail -n1)";
~/.nodenv/bin/nodenv install -s $nodeV;
~/.nodenv/bin/nodenv global $nodeV;
~/.nodenv/shims/npm install --global npm;
~/.nodenv/shims/npm install --global yarn;
~/.nodenv/shims/npm install --global pnpm;
~/.nodenv/bin/nodenv rehash;

# Install Ruby
# ------------------------------------------------------------------------------
if [ -z "$(command -v rbenv)" ]; then
    sudo dnf install -y gcc bzip2 openssl-devel libyaml-devel libffi-devel readline-devel zlib-devel gdbm-devel ncurses-devel;
    rm -rf ~/.rbenv;
    git clone https://github.com/rbenv/rbenv.git ~/.rbenv;
    cd ~/.rbenv && src/configure && make -C src && cd -;
    mkdir -p ~/.rbenv/plugins;
    git clone https://github.com/rbenv/ruby-build.git ~/.rbenv/plugins/ruby-build;
fi
cd ~/.rbenv && git pull && cd -;
cd ~/.rbenv/plugins/ruby-build && git pull && cd -;
rubyV="$(echo "$(~/.rbenv/bin/rbenv install --list)" | awk '{$1=$1};1' | grep '^[0-9]' | grep -v '[a-zA-Z]' | sed '/-/!{s/$/_/}' | sort -V | sed 's/_$//' | tail -n1)";
~/.rbenv/bin/rbenv install -s $rubyV;
~/.rbenv/bin/rbenv global $rubyV;
~/.rbenv/bin/rbenv rehash;

# Install Python
# ------------------------------------------------------------------------------
if [ -z "$(command -v pyenv)" ]; then
    sudo dnf install -y make gcc zlib-devel bzip2 bzip2-devel readline-devel sqlite sqlite-devel openssl-devel tk-devel libffi-devel;
    rm -rf ~/.pyenv;
    git clone https://github.com/pyenv/pyenv.git ~/.pyenv;
    cd ~/.pyenv && src/configure && make -C src && cd -;
fi
cd ~/.pyenv && git pull && cd -;
pythonV="$(echo "$(~/.pyenv/bin/pyenv install --list)" | awk '{$1=$1};1' | grep '^[0-9]' | grep -v '[a-zA-Z]' | sed '/-/!{s/$/_/}' | sort -V | sed 's/_$//' | tail -n1)";
~/.pyenv/bin/pyenv install -s $pythonV;
~/.pyenv/bin/pyenv global $pythonV;
~/.pyenv/bin/pyenv rehash;

# Install Packer
# ------------------------------------------------------------------------------
if [ -z "$(command -v pkenv)" ]; then
    rm -rf ~/.pkenv;
    git clone https://github.com/iamhsa/pkenv.git ~/.pkenv;
fi
cd ~/.pkenv && git pull && cd -;
packerV="$(echo "$(~/.pkenv/bin/pkenv list-remote)" | awk '{$1=$1};1' | grep '^[0-9]' | grep -v '[a-zA-Z]' | sed '/-/!{s/$/_/}' | sort -V | sed 's/_$//' | tail -n1)";
~/.pkenv/bin/pkenv install $packerV;
~/.pkenv/bin/pkenv use $packerV;

# Install Terraform
# ------------------------------------------------------------------------------
if [ -z "$(command -v tfenv)" ]; then
    rm -rf ~/.tfenv;
    git clone https://github.com/tfutils/tfenv.git ~/.tfenv;
fi
cd ~/.tfenv && git pull && cd -;
tfV="$(echo "$(~/.tfenv/bin/tfenv list-remote)" | awk '{$1=$1};1' | grep '^[0-9]' | grep -v '[a-zA-Z]' | sed '/-/!{s/$/_/}' | sort -V | sed 's/_$//' | tail -n1)";
~/.tfenv/bin/tfenv install $tfV;
~/.tfenv/bin/tfenv use $tfV;

# Install Java / Kotlin
# ------------------------------------------------------------------------------
if ! [ -d "$HOME/.sdkman" ]; then
    curl "https://get.sdkman.io?rcupdate=false" | bash;
fi
bash ~/.local/bin/update-sdkman;
