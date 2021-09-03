#!/usr/bin/env bash
set -euo pipefail;

# Source our bashrc file
[ -f ~/.bashrc ] && set +u && . ~/.bashrc && set -u;

# Update the entire system
sudo dnf update -y;

echo ">>> Install Docker (Podman)";
echo "------------------------------------------------------------------------------";
sudo dnf install -y podman;
sudo dnf reinstall -y shadow-utils;
#if ! [ -x "/usr/bin/docker" ]; then
#	sudo ln -s "/usr/bin/podman" "/usr/bin/docker";
#fi
echo "";

echo ">>> Install the ASDF Tool Version Manager";
echo "------------------------------------------------------------------------------";
if [ -z "$(command -v asdf)" ]; then
    rm -rf ~/.asdf;
	sudo dnf install -y curl git;
	git clone https://github.com/asdf-vm/asdf.git ~/.asdf;
	cd ~/.asdf && git checkout "$(git describe --abbrev=0 --tags)" && cd -;
    . "$HOME/.asdf/asdf.sh";
fi
asdf update;
echo "";

echo ">>> Install awscli";
echo "------------------------------------------------------------------------------";
asdf plugin add awscli https://github.com/MetricMike/asdf-awscli.git || true;
asdf plugin update awscli;
asdf install awscli latest;
asdf global awscli latest;
echo "";

echo ">>> Install dotnet";
echo "------------------------------------------------------------------------------";
asdf plugin add dotnet-core https://github.com/emersonsoares/asdf-dotnet-core.git || true;
asdf plugin update dotnet-core;
asdf install dotnet-core latest;
asdf global dotnet-core latest;
echo "";

echo ">>> Install Deno";
echo "------------------------------------------------------------------------------";
asdf plugin add deno https://github.com/asdf-community/asdf-deno.git || true;
asdf plugin update deno;
asdf install deno latest;
asdf global deno latest;
asdf exec deno install -A -f -n udd https://deno.land/x/udd/main.ts;
asdf reshim deno latest;
echo "";

echo ">>> Install Dartlang";
echo "------------------------------------------------------------------------------";
asdf plugin add dart https://github.com/patoconnor43/asdf-dart.git || true;
asdf plugin update dart;
asdf install dart latest;
asdf global dart latest;
echo "";

echo ">>> Install Golang";
echo "------------------------------------------------------------------------------";
asdf plugin add golang https://github.com/kennyp/asdf-golang.git || true;
asdf plugin update golang;
asdf install golang latest;
asdf global golang latest;
asdf exec go install golang.org/x/tools/gopls@latest;
asdf reshim golang latest;
echo "";

echo ">>> Install Nodejs";
echo "------------------------------------------------------------------------------";
asdf plugin add nodejs https://github.com/asdf-vm/asdf-nodejs.git || true;
asdf plugin update nodejs;
asdf install nodejs latest;
asdf global nodejs latest;
asdf exec npm install --global npm;
asdf exec npm install --global yarn;
asdf exec npm install --global pnpm;
echo "";

echo ">>> Install Ruby";
echo "------------------------------------------------------------------------------";
sudo dnf install -y gcc make bzip2 openssl-devel libyaml-devel libffi-devel readline-devel zlib-devel gdbm-devel ncurses-devel
asdf plugin add ruby https://github.com/asdf-vm/asdf-ruby.git || true;
asdf plugin update ruby;
asdf install ruby latest;
asdf global ruby latest;
echo "";

echo ">>> Install Python";
echo "------------------------------------------------------------------------------";
sudo dnf install -y make gcc zlib-devel bzip2 bzip2-devel readline-devel sqlite sqlite-devel openssl-devel tk-devel libffi-devel xz-devel;
asdf plugin add python https://github.com/danhper/asdf-python.git || true;
asdf plugin update python;
asdf install python latest;
asdf global python latest;
echo "";

echo ">>> Install Packer";
echo "------------------------------------------------------------------------------";
asdf plugin add packer https://github.com/asdf-community/asdf-hashicorp.git || true;
asdf plugin update packer;
asdf install packer latest;
asdf global packer latest;
echo "";

echo ">>> Install Terraform";
echo "------------------------------------------------------------------------------";
asdf plugin add terraform https://github.com/asdf-community/asdf-hashicorp.git || true;
asdf plugin update terraform;
asdf install terraform latest;
asdf global terraform latest;
echo "";

echo ">>> Install Java";
echo "------------------------------------------------------------------------------";
asdf plugin add java https://github.com/halcyon/asdf-java.git || true;
asdf plugin update java;
javaV="$(asdf list-all java corretto | grep -v "corretto-musl" | tail -n 1)";
asdf install java ${javaV};
asdf global java ${javaV};
echo "";

echo ">>> Install Maven";
echo "------------------------------------------------------------------------------";
asdf plugin add maven https://github.com/halcyon/asdf-maven.git || true;
asdf plugin update maven;
asdf install maven latest;
asdf global maven latest;
echo "";

echo ">>> Install Kotlin";
echo "------------------------------------------------------------------------------";
asdf plugin add kotlin https://github.com/asdf-community/asdf-kotlin.git || true;
asdf plugin update kotlin;
asdf install kotlin latest;
asdf global kotlin latest;
echo "";

echo ">>> Install Kubectl";
echo "------------------------------------------------------------------------------";
asdf plugin add kubectl https://github.com/asdf-community/asdf-kubectl.git || true;
asdf plugin update kubectl;
asdf install kubectl latest;
asdf global kubectl latest;
echo "";

echo ">>> Install k9s";
echo "------------------------------------------------------------------------------";
asdf plugin add k9s https://github.com/looztra/asdf-k9s.git || true;
asdf plugin update k9s;
asdf install k9s latest;
asdf global k9s latest;
echo "";

echo ">>> Install task";
echo "------------------------------------------------------------------------------";
asdf plugin add task https://github.com/particledecay/asdf-task.git || true;
asdf plugin update task;
asdf install task latest;
asdf global task latest;
echo "";

echo ">>> Install batect";
echo "------------------------------------------------------------------------------";
asdf plugin add batect https://github.com/johnlayton/asdf-batect.git || true;
asdf plugin update batect;
asdf install batect latest;
asdf global batect latest;
echo "";

echo ">>> Install Dprint";
echo "------------------------------------------------------------------------------";
curl -fsSL https://dprint.dev/install.sh | sh;
echo "";
