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
	Exec -ScriptBlock {
		Import-Module PowerShellGet -Force;
		Install-Module PSReadLine -Force;
	}

	# Install the CredentialManager module
	# --------------------------------------------------------------------------
	# see: <https://www.powershellgallery.com/packages/CredentialManager/1.0>
	# 
	# NOTE: We did this hack because it's hard to check to see if the
	# CredentialManager is installed from inside powershell core. It's not
	# impossible I'm just being lazy here. The assumption is if PSReadLine
	# is installed then so too will be CredentialManager because we installed
	# it at the same time.
	Exec -ScriptBlock {
		powershell.exe -Command '& { Import-Module PowerShellGet -Force; Install-Module CredentialManager -Force; }';
	}
}

# Install SSH Agent
# ------------------------------------------------------------------------------
if ((Get-Service -Name ssh-agent).Status -ne "Running") {
	Exec -ScriptBlock { sudo "$env:USERPROFILE\scoop\apps\win32-openssh\current\install-sshd.ps1"; }
	Exec -ScriptBlock { sudo Set-Service -Name ssh-agent -StartupType Automatic; }
	Exec -ScriptBlock { sudo Start-Service -Name ssh-agent; }
}

# Install our run at login script
# ------------------------------------------------------------------------------
if (-Not (Get-ScheduledTask -TaskName "Run at Logon" -ErrorAction Ignore)) {
	$Stt = New-ScheduledTaskTrigger -AtLogOn;
	
	$Sta = New-ScheduledTaskAction `
		-Execute "C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe" `
		-Argument "-NoLogo -NoProfile -File .\run-at-logon.ps1" `
		-WorkingDirectory "$env:USERPROFILE\Documents\WindowsPowershell\Scripts";
	
	$u = whoami;
	$STPrincipal = New-ScheduledTaskPrincipal -UserID "$u";
	
	Register-ScheduledTask "Run at Logon" `
		-Principal $STPrincipal `
		-Trigger $Stt `
		-Action $Sta;
}

# Install the rest of our apps
# ------------------------------------------------------------------------------
scoop install `
	adoptopenjdk-hotspot`
	aws `
	aws-vault `
	dart `
	dotnet-sdk `
	drun `
	gitkraken `
	go `
	kotlin `
	ktlint `
	nuget `
	maven `
	nodejs `
	nssm `
	packer `
	python `
	ruby `
	terraform `
	vscode `
	wavebox10-pro `
	yarn;

# Install additional SSH Keys
# ------------------------------------------------------------------------------
# TODO: Would be nice use gopass as an actual ssh-agent?
RmIfExists -Path $env:USERPROFILE/.ssh/keys;

mkdir $env:USERPROFILE/.ssh/keys/xero-payroll-prod;
gopass bin cp keys/ssh/xero-payroll-prod/payroll-checkpoint.pem $env:USERPROFILE/.ssh/keys/xero-payroll-prod/payroll-checkpoint.pem;
gopass bin cp keys/ssh/xero-payroll-prod/payroll-dev-public.pem $env:USERPROFILE/.ssh/keys/xero-payroll-prod/payroll-dev-public.pem;
gopass bin cp keys/ssh/xero-payroll-prod/payroll-devops.pem $env:USERPROFILE/.ssh/keys/xero-payroll-prod/payroll-devops.pem;

mkdir $env:USERPROFILE/.ssh/keys/xero-payroll-test;
gopass bin cp keys/ssh/xero-payroll-test/payroll-checkpoint.pem $env:USERPROFILE/.ssh/keys/xero-payroll-test/payroll-checkpoint.pem;
gopass bin cp keys/ssh/xero-payroll-test/payroll-dev-public.pem $env:USERPROFILE/.ssh/keys/xero-payroll-test/payroll-dev-public.pem;
gopass bin cp keys/ssh/xero-payroll-test/payroll-devops.pem $env:USERPROFILE/.ssh/keys/xero-payroll-test/payroll-devops.pem;

mkdir $env:USERPROFILE/.ssh/keys/xero-payroll-uat;
gopass bin cp keys/ssh/xero-payroll-uat/payroll-checkpoint.pem $env:USERPROFILE/.ssh/keys/xero-payroll-uat/payroll-checkpoint.pem;
gopass bin cp keys/ssh/xero-payroll-uat/payroll-dev-public.pem $env:USERPROFILE/.ssh/keys/xero-payroll-uat/payroll-dev-public.pem;
gopass bin cp keys/ssh/xero-payroll-uat/payroll-devops.pem $env:USERPROFILE/.ssh/keys/xero-payroll-uat/payroll-devops.pem;

mkdir $env:USERPROFILE/.ssh/keys/xero-ps-paas-svc;
gopass bin cp keys/ssh/xero-ps-paas-svc/payroll-devops.pem $env:USERPROFILE/.ssh/keys/xero-ps-paas-svc/payroll-devops.pem;