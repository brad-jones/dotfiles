# Make sure this script dies when something bad happens
# A really crappy version of bash set -e;
$ErrorActionPreference = 'stop';
function Exec {
    param (
        [scriptblock]$ScriptBlock,
        [string]$ErrorAction = $ErrorActionPreference
    )
    & @ScriptBlock
    if (($lastexitcode -ne 0) -and $ErrorAction -eq "Stop") {
        exit $lastexitcode
    }
}

# Remove a path quietly
Function RmIfExists {
	Param ($Path);
	if (Test-Path -Path $Path) {
		Remove-Item -Path $Path -Recurse -Force;
	}
}

# Check if a command exists or not
Function CommandExists {
	Param ($Command);
	$oldPreference = $ErrorActionPreference;
	$ErrorActionPreference = 'stop';
	try {
		if (Get-Command $Command) {
			return $true;
		}
	} catch {
		return $false;
	} finally {
		$ErrorActionPreference = $oldPreference;
	}
}

# Ensure scoop is installed
if (!(CommandExists -Command "scoop")) {
	iwr -useb get.scoop.sh | iex;
}

# Ensure powershell core is installed and then re-execute using it
if ($PSVersionTable.PSVersion.Major -le 5) {
	if (!(CommandExists -Command "pwsh")) {=
		scoop install pwsh;
	}
	pwsh $PSCommandPath;
	exit;
}

# Install some tools
Exec -ScriptBlock { scoop install sudo 7zip jq wget grep sed git gpg win32-openssh gopass }

# Install chezmoi
RmIfExists -Path $env:TEMP\chezmoi;
RmIfExists -Path $env:TEMP\chezmoi.zip;
RmIfExists -Path $env:USERPROFILE\.local\bin\chezmoi.exe;
$ErrorActionPreference = 'continue';
$chezmoiV = wget https://github.com/twpayne/chezmoi/releases/latest -O /dev/null 2>&1 | grep Location: | sed -r 's~^.*tag/v(.*?) \[.*~\1~g';
$ErrorActionPreference = 'stop';
Exec -ScriptBlock { wget https://github.com/twpayne/chezmoi/releases/download/v${chezmoiV}/chezmoi_${chezmoiV}_windows_amd64.zip -O $env:TEMP\chezmoi.zip; }
Exec -ScriptBlock { 7z x $env:TEMP\chezmoi.zip "-o${env:TEMP}\chezmoi"; }
Copy-Item -Path $env:TEMP\chezmoi\chezmoi.exe -Destination $env:USERPROFILE\.local\bin\chezmoi.exe;
RmIfExists -Path $env:TEMP\chezmoi.zip; RmIfExists -Path $env:TEMP\chezmoi;
if (!($env:PATH -like "*$env:UserProfile\.local\bin*")) {
	$env:PATH += ";$env:UserProfile\.local\bin";
	[Environment]::SetEnvironmentVariable("PATH", $env:PATH, "User");
}

# Ensure this script is Idempotent
RmIfExists -Path "$env:TEMP\vault-key";
RmIfExists -Path "$env:TEMP\brad@bjc.id.au";
RmIfExists -Path "$env:TEMP\brad.jones@xero.com";
RmIfExists -Path "$env:USERPROFILE\.password-store";
RmIfExists -Path "$env:USERPROFILE\.local\share\chezmoi";
RmIfExists -Path "$env:USERPROFILE\.ssh\brad@bjc.id.au";
RmIfExists -Path "$env:USERPROFILE\.ssh\brad@bjc.id.au.pub";
RmIfExists -Path "$env:USERPROFILE\.ssh\brad.jones@xero.com";
RmIfExists -Path "$env:USERPROFILE\.ssh\brad.jones@xero.com.pub";

# Ensure git doesn't do silly things with Windows line endings
Exec -ScriptBlock { git config --global core.eol lf; }
Exec -ScriptBlock { git config --global core.autocrlf false; }

# Install the GPG key from gitlab that is used to decrypt my gopass vault
Exec -ScriptBlock { git clone https://gitlab.com/brad-jones/vault-key.git "$env:TEMP\vault-key"; }
Exec -ScriptBlock { gpg --import "$env:TEMP\vault-key\private.pem"; }
Exec -ScriptBlock { echo "5`r`ny" | gpg --command-fd 0 --edit-key "Brad Jones (vault) <brad@bjc.id.au>" trust; }
RmIfExists -Path "$env:TEMP\vault-key";

# Install the gopass vault from github. To unlock the vault we need to know
# 3 things (gitlab password, github password & the key passphrase).
Exec -ScriptBlock { git clone https://github.com/brad-jones/vault.git "$env:USERPROFILE\.password-store"; }
Exec -ScriptBlock { git --git-dir "$env:USERPROFILE\.password-store\.git" remote set-url origin git@github.com:brad-jones/vault.git; }

# Install my personal and work SSH keys
Exec -ScriptBlock { gopass bin cp "keys/ssh/brad@bjc.id.au" "$env:USERPROFILE\.ssh\brad@bjc.id.au"; }
Exec -ScriptBlock { gopass bin cp "keys/ssh/brad.jones@xero.com" "$env:USERPROFILE\.ssh\brad.jones@xero.com"; }

# Install my personal GPG key
Exec -ScriptBlock { gopass bin cp "keys/gpg/brad@bjc.id.au" "$env:TEMP/brad@bjc.id.au"; }
Exec -ScriptBlock { gpg --import "$env:TEMP/brad@bjc.id.au"; }
Exec -ScriptBlock { echo "5`r`ny" | gpg --command-fd 0 --edit-key "Brad Jones <brad@bjc.id.au>" trust; }
RmIfExists -Path "$env:TEMP/brad@bjc.id.au";

# Install my work GPG key
Exec -ScriptBlock { gopass bin cp "keys/gpg/brad.jones@xero.com" "$env:TEMP/brad.jones@xero.com"; }
Exec -ScriptBlock { gpg --import "$env:TEMP/brad.jones@xero.com"; }
Exec -ScriptBlock { echo "5`r`ny" | gpg --command-fd 0 --edit-key "Brad Jones <brad.jones@xero.com>" trust; }
RmIfExists -Path "$env:TEMP/brad.jones@xero.com";

# Install my dotfiles
Exec -ScriptBlock { chezmoi init https://github.com/brad-jones/dotfiles.git; }
Exec -ScriptBlock { cmdkey /delete:"git:https://github.com"; }
Exec -ScriptBlock { cmdkey /delete:"git:https://gitlab.com"; }
Exec -ScriptBlock { cmdkey /delete:"git:https://brad@bjc.id.au@gitlab.com"; }
$gDir = chezmoi source-path;
Exec -ScriptBlock { git --git-dir "$gDir\.git" remote set-url origin git@github.com:brad-jones/dotfiles.git; }
Exec -ScriptBlock { chezmoi apply --debug; }

# Reboot to make sure things like kernels are updated etc
#Restart-Computer;