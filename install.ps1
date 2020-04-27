Function RmIfExists {
	Param ($Path);
	if (Test-Path -Path $Path) {
		Remove-Item -Path $Path -Recurse -Force;
	}
}

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
		$ErrorActionPreference=$oldPreference;
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
scoop install 7zip jq wget grep sed git gpg win32-openssh gopass;

# Install chezmoi
RmIfExists -Path $env:TEMP\chezmoi;
RmIfExists -Path $env:TEMP\chezmoi.zip;
RmIfExists -Path $env:USERPROFILE\.local\bin\chezmoi.exe;
$chezmoiV = wget https://github.com/twpayne/chezmoi/releases/latest -O /dev/null 2>&1 | grep Location: | sed -r 's~^.*tag/v(.*?) \[.*~\1~g';
wget https://github.com/twpayne/chezmoi/releases/download/v${chezmoiV}/chezmoi_${chezmoiV}_windows_amd64.zip -O $env:TEMP\chezmoi.zip;
7z x $env:TEMP\chezmoi.zip "-o${env:TEMP}\chezmoi";
Copy-Item -Path $env:TEMP\chezmoi\chezmoi.exe -Destination $env:USERPROFILE\.local\bin\chezmoi.exe;
RmIfExists -Path $env:TEMP\chezmoi.zip; RmIfExists -Path $env:TEMP\chezmoi;
if (!($env:PATH -like "*$env:UserProfile\.local\bin*")) {
	$env:PATH += ";$env:UserProfile\.local\bin";
	[Environment]::SetEnvironmentVariable("PATH", $env:PATH, "User");
}

RmIfExists -Path "$env:TEMP\vault-key";
RmIfExists -Path "$env:TEMP\brad@bjc.id.au";
RmIfExists -Path "$env:TEMP\brad.jones@xero.com";
RmIfExists -Path "$env:USERPROFILE\.password-store";
RmIfExists -Path "$env:USERPROFILE\.local\share\chezmoi";
RmIfExists -Path "$env:USERPROFILE\.ssh\brad@bjc.id.au";
RmIfExists -Path "$env:USERPROFILE\.ssh\brad@bjc.id.au.pub";
RmIfExists -Path "$env:USERPROFILE\.ssh\brad.jones@xero.com";
RmIfExists -Path "$env:USERPROFILE\.ssh\brad.jones@xero.com.pub";

git clone https://gitlab.com/brad-jones/vault-key.git "$env:TEMP\vault-key";
gpg --import "$env:TEMP\vault-key\private.pem";
echo "5`r`ny" | gpg --command-fd 0 --edit-key "Brad Jones (vault) <brad@bjc.id.au>" trust;
RmIfExists -Path "$env:TEMP\vault-key";

git clone https://github.com/brad-jones/vault.git "$env:USERPROFILE\.password-store";
git --git-dir "$env:USERPROFILE\.password-store\.git" remote set-url origin git@github.com:brad-jones/vault.git;
gopass bin cp "keys/ssh/brad@bjc.id.au" "$env:USERPROFILE\.ssh\brad@bjc.id.au";
gopass bin cp "keys/ssh/brad.jones@xero.com" "$env:USERPROFILE\.ssh\brad.jones@xero.com";

gopass bin cp "keys/gpg/brad@bjc.id.au" "$env:TEMP/brad@bjc.id.au";
gpg --import "$env:TEMP/brad@bjc.id.au";
echo "5`r`ny" | gpg --command-fd 0 --edit-key "Brad Jones <brad@bjc.id.au>" trust;
RmIfExists -Path "$env:TEMP/brad@bjc.id.au";

gopass bin cp "keys/gpg/brad.jones@xero.com" "$env:TEMP/brad.jones@xero.com";
gpg --import "$env:TEMP/brad.jones@xero.com";
echo "5`r`ny" | gpg --command-fd 0 --edit-key "Brad Jones <brad.jones@xero.com>" trust;
RmIfExists -Path "$env:TEMP/brad.jones@xero.com";

chezmoi init https://github.com/brad-jones/dotfiles.git;
$gDir = chezmoi source-path;
git --git-dir "$gDir\.git" remote set-url origin git@github.com:brad-jones/dotfiles.git;
chezmoi apply --debug;
