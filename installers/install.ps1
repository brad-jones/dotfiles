#!/usr/bin/env pwsh
$ErrorActionPreference = 'Stop';

# ------------------------------------------------------------------------------
# CONFIG
# ------------------------------------------------------------------------------

# The version of the executable to download & install
$version = $env:DOTFILES_VERSION;
$version = if ([string]::IsNullOrWhiteSpace($version)) {
	"0.0.0"
}

# Where to install the executable, the parent directory.
$installDir = $env:DOTFILES_INSTALL_DIR;
$installDir = if ([string]::IsNullOrWhiteSpace($installDir)) {
	"$env:USERPROFILE\.local\bin"
}

# The final url that will be downloaded.
$finalUrl = $env:DOTFILES_DOWNLOAD_URL;
$finalUrl = if ([string]::IsNullOrWhiteSpace($finalUrl)) {
	"https://github.com/brad-jones/dotfiles/releases/${version}/download/dotfiles_windows_amd64.exe"
}

# The final output path of the downloaded executable.
$finalOutput = $env:DOTFILES_OUTPUT;
$finalOutput = if ([string]::IsNullOrWhiteSpace($finalOutput)) {
	"$installDir\dotfiles.exe"
}

# Comma-separated list of codenotary.io SignerID(s)
$signers = $env:DOTFILES_SIGNERS;
$signers = if ([string]::IsNullOrWhiteSpace($signers)) {
	"0x33dc7b56c2be6175e1ad17e31f003f55943fa4ce"
}

# By default this script will execute the downloaded & authenticated executable,
# if for some reason you do not want to happen you can set this value to false.
$runAfterInstall = $env:DOTFILES_RUN_AFTER_INSTALL;
$runAfterInstall = if ([string]::IsNullOrWhiteSpace($runAfterInstall)) {
	"true"
}

# ------------------------------------------------------------------------------
# HELPER FUNCTIONS
# ------------------------------------------------------------------------------

function CommandExists {
	Param ($Command);
	try {
		if (Get-Command $Command) {
			return $true;
		}
	} catch {
		return $false;
	}
}

# ------------------------------------------------------------------------------
# MAIN ENTRYPOINT
# ------------------------------------------------------------------------------

# Make sure the finalOutput dir exists
$parentDir = Split-Path $finalOutput;
[Console]::Write("Ensuring $parentDir exists... ");
New-Item -ItemType "directory" -Path "$parentDir" -Force > $null;
[Console]::WriteLine("DONE");

# Download the executable.
#
# The normal "Invoke-WebRequest" is super-duper slow in PowerShell Desktop,
# so we will attempt to use BitsTransfer instead, if available.
#
# In PowerShell 7+ they seem to have resolved this performance issue.
[Console]::Write("Downloading $finalUrl to $finalOutput... ");
if ($PSVersionTable.PSVersion.Major -eq 5 -and $(CommandExists "Start-BitsTransfer")) {
	Start-BitsTransfer -Source "$finalUrl" -Destination "$finalOutput";
} else {
	Invoke-WebRequest -UseBasicParsing "$finalUrl" -OutFile "$finalOutput";
}
[Console]::WriteLine("DONE");

# Authenticate the downloaded executable against codenotary.io
[Console]::Write("Authenticating $finalOutput against codenotary.io... ");
$hash = (Get-FileHash "$finalOutput" -Algorithm "SHA256").Hash;
$results = Invoke-RestMethod "https://api.codenotary.io/authenticate/${hash}?signers=${$signers}";
if ($results.verification.status -ne 0) {
	[Console]::ForegroundColor = 'red';
	[Console]::Error.WriteLine("UNTRUSTED");
	[Console]::Error.WriteLine("!!! Downloaded executable does not authenticate against codenotary.io !!!");
	Remove-Item -Path "$finalOutput";
	[Console]::Error.WriteLine("$finalOutput has been deleted");
	[Console]::ResetColor();
	exit $results.verification.status;
}
[Console]::WriteLine("TRUSTED");

# Execute the executable
if ($runAfterInstall -eq "true") {
	[Console]::WriteLine("Executing $finalOutput...");
	& "$finalOutput" $args;
	if ($LASTEXITCODE) {
		exit $LASTEXITCODE;
	}
}
