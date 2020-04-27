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

# Remove Linux specfic stuff
# ------------------------------------------------------------------------------
RmIfExists -Path ~/.config/systemd;

# Install a semi-sensible readline for powershell
# ------------------------------------------------------------------------------
# see: <https://github.com/PowerShell/PSReadLine>
Exec -ScriptBlock { Import-Module PowerShellGet -Force; Install-Module PSReadLine -Force; }
powershell.exe -Command '& { Import-Module PowerShellGet -Force; Install-Module PSReadLine -Force; }';