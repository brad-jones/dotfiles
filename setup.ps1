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

Function ModuleExists {
	Param ($Module);
	$oldPreference = $ErrorActionPreference;
	$ErrorActionPreference = 'stop';
	try {
		if (Get-InstalledModule -Name $Module) {
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
if (-Not (ModuleExists -Module PSReadLine)) {
	Exec -ScriptBlock { Import-Module PowerShellGet -Force; Install-Module PSReadLine -Force; }
	powershell.exe -Command '& { Import-Module PowerShellGet -Force; Install-Module PSReadLine -Force; }';
}

# Install SSH Agent
# ------------------------------------------------------------------------------
if ((Get-Service -Name ssh-agent).Status -ne "Running") {
	Exec -ScriptBlock { sudo "$env:USERPROFILE\scoop\apps\win32-openssh\current\install-sshd.ps1"; }
	Exec -ScriptBlock { sudo Set-Service -Name ssh-agent -StartupType Automatic; }
	Exec -ScriptBlock { sudo Start-Service -Name ssh-agent; }
}

# Install Docker
# ------------------------------------------------------------------------------

# Install Dartlang
# ------------------------------------------------------------------------------

# Install Dotnet Core
# ------------------------------------------------------------------------------

# Install Golang
# ------------------------------------------------------------------------------

# Install Nodejs
# ------------------------------------------------------------------------------

# Install Ruby
# ------------------------------------------------------------------------------

# Install Python
# ------------------------------------------------------------------------------

# Install awscli
# ------------------------------------------------------------------------------

# Install aws-vault
# ------------------------------------------------------------------------------

# Install Packer
# ------------------------------------------------------------------------------

# Install Terraform
# ------------------------------------------------------------------------------

# Install Java / Kotlin
# ------------------------------------------------------------------------------

# Install additional SSH Keys
# ------------------------------------------------------------------------------
# TODO: Would be nice use gopass as an actual ssh-agent?
RmIfExists -Path ~/.ssh/keys;

mkdir ~/.ssh/keys/xero-payroll-prod;
gopass bin cp keys/ssh/xero-payroll-prod/payroll-checkpoint.pem ~/.ssh/keys/xero-payroll-prod/payroll-checkpoint.pem;
gopass bin cp keys/ssh/xero-payroll-prod/payroll-dev-public.pem ~/.ssh/keys/xero-payroll-prod/payroll-dev-public.pem;
gopass bin cp keys/ssh/xero-payroll-prod/payroll-devops.pem ~/.ssh/keys/xero-payroll-prod/payroll-devops.pem;

mkdir ~/.ssh/keys/xero-payroll-test;
gopass bin cp keys/ssh/xero-payroll-test/payroll-checkpoint.pem ~/.ssh/keys/xero-payroll-test/payroll-checkpoint.pem;
gopass bin cp keys/ssh/xero-payroll-test/payroll-dev-public.pem ~/.ssh/keys/xero-payroll-test/payroll-dev-public.pem;
gopass bin cp keys/ssh/xero-payroll-test/payroll-devops.pem ~/.ssh/keys/xero-payroll-test/payroll-devops.pem;

mkdir ~/.ssh/keys/xero-payroll-uat;
gopass bin cp keys/ssh/xero-payroll-uat/payroll-checkpoint.pem ~/.ssh/keys/xero-payroll-uat/payroll-checkpoint.pem;
gopass bin cp keys/ssh/xero-payroll-uat/payroll-dev-public.pem ~/.ssh/keys/xero-payroll-uat/payroll-dev-public.pem;
gopass bin cp keys/ssh/xero-payroll-uat/payroll-devops.pem ~/.ssh/keys/xero-payroll-uat/payroll-devops.pem;

mkdir ~/.ssh/keys/xero-ps-paas-svc;
gopass bin cp keys/ssh/xero-ps-paas-svc/payroll-devops.pem ~/.ssh/keys/xero-ps-paas-svc/payroll-devops.pem;