#!/usr/bin/env bash
set -euxo pipefail;

# Remove Windows specfic stuff
# ------------------------------------------------------------------------------
rm -rf ~/Documents;

# Install SSH/GPG Agent
# ------------------------------------------------------------------------------
# On desktop systems, ssh/gpg-agent is already taken care of for us.
# We only need to set this up on headless systems.
if [ -z ${XDG_SESSION_DESKTOP+x} ];
then
	sudo loginctl enable-linger $USER;
	systemctl --user daemon-reload;
	systemctl --user enable gpg-agent.socket --now;
    systemctl --user enable ssh-agent.service --now;
fi

# Install Homebrew
# ------------------------------------------------------------------------------
if ! [ -x "$(command -v brew)" ]; then
    rm -rf ~/.linuxbrew;
    sudo dnf install -y curl file git libxcrypt-compat;
    git clone https://github.com/Homebrew/brew ~/.linuxbrew/Homebrew;
    mkdir ~/.linuxbrew/bin;
    ln -s ~/.linuxbrew/Homebrew/bin/brew ~/.linuxbrew/bin;
    eval "$(~/.linuxbrew/bin/brew shellenv)";
    ~/.linuxbrew/bin/brew --version;
fi

# Install Docker
# ------------------------------------------------------------------------------
if ! [ -x "$(command -v docker)" ]; then
    sudo dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo;
    sudo dnf install -y docker-ce docker-ce-cli containerd.io grubby;
    sudo grubby --update-kernel=ALL --args="systemd.unified_cgroup_hierarchy=0";
    sudo systemctl enable docker;
    sudo usermod -aG docker $USER;

    # When running docker inside a hyper-v vm sometimes we get IP clashes because
    # the default switch for hyper-v also likes to make use of the 172.x range.
    # see: https://stackoverflow.com/questions/44003663
    # also: https://github.com/moby/moby/pull/29376
    dockerConfig=$(cat <<'EOF'
{
    "experimental": true,
    "default-address-pools": [
        {
            "base": "10.10.0.0/16",
            "size": 24
        }
    ]
}
EOF
)

    sudo sh -c "echo '$dockerConfig' > /etc/docker/daemon.json";
fi

# Install Dartlang
# ------------------------------------------------------------------------------
if ! [ -x "$(command -v dart)" ]; then
    rm -rf ~/.dart;
    mkdir ~/.dart;
    tmpFolder="/tmp/$(uuidgen)";
    mkdir -p $tmpFolder;
    function finish {
        rm -rf $tmpFolder;
    }
    curl "https://storage.googleapis.com/dart-archive/channels/stable/release/latest/sdk/dartsdk-linux-x64-release.zip" -o "$tmpFolder/dartsdk.zip";
    unzip "$tmpFolder/dartsdk.zip" -d "$tmpFolder/extracted";
    dartV="$(cat "$tmpFolder/extracted/dart-sdk/version")";
    mv "$tmpFolder/extracted/dart-sdk" "$HOME/.dart/$dartV";
    ln -s "$HOME/.dart/$dartV" "$HOME/.dart/current";
    cd ~/.local/sbin && ~/.dart/current/bin/pub get && cd -;
fi

# Install Dotnet Core
# ------------------------------------------------------------------------------
if ! [ -x "$(command -v dotnet)" ]; then
    rm -rf ~/.dotnet;
    sudo rpm --import "https://keyserver.ubuntu.com/pks/lookup?op=get&search=0x3FA7E0328081BFF6A14DA29AA6A19B38D3D831EF";
    sudo dnf config-manager --add-repo https://download.mono-project.com/repo/centos8-stable.repo;
    sudo dnf install -y krb5-libs libcurl libgdiplus libicu libunwind libuuid lttng-ust openssl-libs zlib;
    curl https://dotnet.microsoft.com/download/dotnet-core/scripts/v1/dotnet-install.sh | bash;
    ~/.dotnet/dotnet --info;
fi

# Install Powershell Core
# ------------------------------------------------------------------------------
if ! [ -x "$(command -v pwsh)" ]; then
    sudo rpm --import https://packages.microsoft.com/keys/microsoft.asc;
    sudo dnf config-manager --add-repo https://packages.microsoft.com/config/rhel/7/prod.repo
    sudo dnf install -y compat-openssl10;
    sudo dnf install -y powershell;
fi

# Install Golang
# ------------------------------------------------------------------------------
if ! [ -x "$(command -v goenv)" ]; then
    rm -rf ~/.goenv;
    git clone https://github.com/syndbg/goenv.git ~/.goenv;
    cd ~/.goenv && src/configure && make -C src && cd -;
    goV="$(echo "$(~/.goenv/bin/goenv install --list)" | awk '{$1=$1};1' | grep '^[0-9]' | grep -v '[a-zA-Z]' | sed '/-/!{s/$/_/}' | sort -V | sed 's/_$//' | tail -n1)";
    ~/.goenv/bin/goenv install $goV;
    ~/.goenv/bin/goenv global $goV;
    ~/.goenv/bin/goenv rehash;
fi

# Install Nodejs
# ------------------------------------------------------------------------------
if ! [ -x "$(command -v nodenv)" ]; then
    rm -rf ~/.nodenv;
    git clone https://github.com/nodenv/nodenv.git ~/.nodenv;
    mkdir -p ~/.nodenv/plugins;
    git clone https://github.com/nodenv/node-build.git ~/.nodenv/plugins/node-build;
    git clone https://github.com/nodenv/nodenv-update.git ~/.nodenv/plugins/nodenv-update;
    git clone https://github.com/nodenv/node-build-update-defs.git ~/.nodenv/plugins/node-build-update-defs;
    cd ~/.nodenv && src/configure && make -C src && cd -;
    nodeV="$(echo "$(~/.nodenv/bin/nodenv install --list)" | awk '{$1=$1};1' | grep '^[0-9]' | grep -v '[a-zA-Z]' | sed '/-/!{s/$/_/}' | sort -V | sed 's/_$//' | tail -n1)";
    ~/.nodenv/bin/nodenv install $nodeV;
    ~/.nodenv/bin/nodenv global $nodeV;
    ~/.nodenv/shims/npm install --global yarn;
    ~/.nodenv/shims/npm install --global pnpm;
    ~/.nodenv/bin/nodenv rehash;
fi

# Install Ruby
# ------------------------------------------------------------------------------
if ! [ -x "$(command -v rbenv)" ]; then
    sudo dnf install -y gcc bzip2 openssl-devel libyaml-devel libffi-devel readline-devel zlib-devel gdbm-devel ncurses-devel;
    rm -rf ~/.rbenv;
    git clone https://github.com/rbenv/rbenv.git ~/.rbenv;
    cd ~/.rbenv && src/configure && make -C src && cd -;
    mkdir -p ~/.rbenv/plugins;
    git clone https://github.com/rbenv/ruby-build.git ~/.rbenv/plugins/ruby-build;
    rubyV="$(echo "$(~/.rbenv/bin/rbenv install --list)" | awk '{$1=$1};1' | grep '^[0-9]' | grep -v '[a-zA-Z]' | sed '/-/!{s/$/_/}' | sort -V | sed 's/_$//' | tail -n1)";
    ~/.rbenv/bin/rbenv install $rubyV;
    ~/.rbenv/bin/rbenv global $rubyV;
    ~/.rbenv/bin/rbenv rehash;
fi

# Install Python
# ------------------------------------------------------------------------------
if ! [ -x "$(command -v pyenv)" ]; then
    sudo dnf install -y make gcc zlib-devel bzip2 bzip2-devel readline-devel sqlite sqlite-devel openssl-devel tk-devel libffi-devel;
    rm -rf ~/.pyenv;
    git clone https://github.com/pyenv/pyenv.git ~/.pyenv;
    cd ~/.pyenv && src/configure && make -C src && cd -;
    pythonV="$(echo "$(~/.pyenv/bin/pyenv install --list)" | awk '{$1=$1};1' | grep '^[0-9]' | grep -v '[a-zA-Z]' | sed '/-/!{s/$/_/}' | sort -V | sed 's/_$//' | tail -n1)";
    ~/.pyenv/bin/pyenv install $pythonV;
    ~/.pyenv/bin/pyenv global $pythonV;
    ~/.pyenv/bin/pyenv rehash;
fi

# Install awscli
# ------------------------------------------------------------------------------
if ! [ -x "$(command -v aws)" ]; then
    rm -rf ~/.local/aws-cli ~/.local/bin/aws;
    tmpFolder="/tmp/$(uuidgen)";
    mkdir -p $tmpFolder;
    function finish {
        rm -rf $tmpFolder;
    }
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "$tmpFolder/awscliv2.zip";
    unzip "$tmpFolder/awscliv2.zip" -d "$tmpFolder/extracted";
    $tmpFolder/extracted/aws/install -i ~/.local/aws-cli -b ~/.local/bin;
    ~/.local/bin/aws --version;
fi

# Install aws-vault
# ------------------------------------------------------------------------------
if ! [ -x "$(command -v aws-vault)" ]; then
    ~/.linuxbrew/bin/brew install aws-vault;
fi

# Install Packer
# ------------------------------------------------------------------------------
if ! [ -x "$(command -v pkenv)" ]; then
    rm -rf ~/.pkenv;
    git clone https://github.com/iamhsa/pkenv.git ~/.pkenv;
    packerV="$(echo "$(~/.pkenv/bin/pkenv list-remote)" | awk '{$1=$1};1' | grep '^[0-9]' | grep -v '[a-zA-Z]' | sed '/-/!{s/$/_/}' | sort -V | sed 's/_$//' | tail -n1)";
    ~/.pkenv/bin/pkenv install $packerV;
    ~/.pkenv/bin/pkenv use $packerV;
fi

# Install Terraform
# ------------------------------------------------------------------------------
if ! [ -x "$(command -v tfenv)" ]; then
    rm -rf ~/.tfenv;
    git clone https://github.com/tfutils/tfenv.git ~/.tfenv;
    tfV="$(echo "$(~/.tfenv/bin/tfenv list-remote)" | awk '{$1=$1};1' | grep '^[0-9]' | grep -v '[a-zA-Z]' | sed '/-/!{s/$/_/}' | sort -V | sed 's/_$//' | tail -n1)";
    ~/.tfenv/bin/tfenv install $tfV;
    ~/.tfenv/bin/tfenv use $tfV;
fi


# Install Java / Kotlin
# ------------------------------------------------------------------------------
if ! [ -d ~/.sdkman ]; then
    curl "https://get.sdkman.io?rcupdate=false" | bash;
    bash ~/.local/bin/update-sdkman;
fi

# Install additional SSH Keys
# ------------------------------------------------------------------------------
# TODO: Would be nice use gopass as an actual ssh-agent?
rm -rf ~/.ssh/keys;

mkdir -p ~/.ssh/keys/xero-payroll-prod;
gopass bin cp keys/ssh/xero-payroll-prod/payroll-checkpoint.pem ~/.ssh/keys/xero-payroll-prod/payroll-checkpoint.pem;
gopass bin cp keys/ssh/xero-payroll-prod/payroll-dev-public.pem ~/.ssh/keys/xero-payroll-prod/payroll-dev-public.pem;
gopass bin cp keys/ssh/xero-payroll-prod/payroll-devops.pem ~/.ssh/keys/xero-payroll-prod/payroll-devops.pem;

mkdir -p ~/.ssh/keys/xero-payroll-test;
gopass bin cp keys/ssh/xero-payroll-test/payroll-checkpoint.pem ~/.ssh/keys/xero-payroll-test/payroll-checkpoint.pem;
gopass bin cp keys/ssh/xero-payroll-test/payroll-dev-public.pem ~/.ssh/keys/xero-payroll-test/payroll-dev-public.pem;
gopass bin cp keys/ssh/xero-payroll-test/payroll-devops.pem ~/.ssh/keys/xero-payroll-test/payroll-devops.pem;

mkdir -p ~/.ssh/keys/xero-payroll-uat;
gopass bin cp keys/ssh/xero-payroll-uat/payroll-checkpoint.pem ~/.ssh/keys/xero-payroll-uat/payroll-checkpoint.pem;
gopass bin cp keys/ssh/xero-payroll-uat/payroll-dev-public.pem ~/.ssh/keys/xero-payroll-uat/payroll-dev-public.pem;
gopass bin cp keys/ssh/xero-payroll-uat/payroll-devops.pem ~/.ssh/keys/xero-payroll-uat/payroll-devops.pem;

mkdir -p ~/.ssh/keys/xero-ps-paas-svc;
gopass bin cp keys/ssh/xero-ps-paas-svc/payroll-devops.pem ~/.ssh/keys/xero-ps-paas-svc/payroll-devops.pem;