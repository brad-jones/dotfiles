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

# Retry a given script block until success
function RetryCommand {
	[CmdletBinding()]
	param (
		[parameter(Mandatory, ValueFromPipeline)] 
		[ValidateNotNullOrEmpty()]
		[scriptblock] $ScriptBlock,
		[int] $RetryCount = 5,
		[int] $TimeoutInSecs = 1,
		[string] $SuccessMessage = "Command executed successfuly!",
		[string] $FailureMessage = "Failed to execute the command"
	)
		
	process {
		$Attempt = 1
		$Flag = $true
		
		do {
			try {
				$PreviousPreference = $ErrorActionPreference
				$ErrorActionPreference = 'Stop'
				Invoke-Command -ScriptBlock $ScriptBlock -OutVariable Result              
				$ErrorActionPreference = $PreviousPreference
				Write-Verbose "$SuccessMessage `n"
				$Flag = $false
			}
			catch {
				if ($Attempt -gt $RetryCount) {
					Write-Verbose "$FailureMessage! Total retry attempts: $RetryCount"
					Write-Verbose "[Error Message] $($_.exception.message) `n"
					$Flag = $false
				}
				else {
					Write-Verbose "[$Attempt/$RetryCount] $FailureMessage. Retrying in $TimeoutInSecs seconds..."
					Start-Sleep -Seconds $TimeoutInSecs
					$Attempt = $Attempt + 1
				}
			}
		}
		While ($Flag)
		
	}
}

# Get our passphrases from the windows cred store
echo "Getting passphrases from windows cred store";
$passwordPersonal = Get-StoredCredential -Target "passphrase:brad@bjc.id.au";
$unsecurePasswordPersonal = [System.Net.NetworkCredential]::new('', $passwordPersonal.Password).Password;
$passwordProfessional = Get-StoredCredential -Target "passphrase:brad.jones@xero.com";
$unsecurePasswordProfessional = [System.Net.NetworkCredential]::new('', $passwordProfessional.Password).Password;

# Make sure the windows GPG agent is running
echo "Start the GPG agent";
RetryCommand -Verbose -ScriptBlock {
	gpg-connect-agent /bye;
	if ($LastExitCode -ne 0) {
		throw "failed";
	}
}

# Add the GPG keys to the windows agent
Exec -ScriptBlock { gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "83D182028C7F2DF102F09E61FF308BBB10F539D8"; }
echo "Unlocked: 83D182028C7F2DF102F09E61FF308BBB10F539D8 @ localhost";
Exec -ScriptBlock { gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB"; }
echo "Unlocked: F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB @ localhost";
Exec -ScriptBlock { gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "1A8059A4CC0F06F670492ABBD0053F0772B75829"; }
echo "Unlocked: 1A8059A4CC0F06F670492ABBD0053F0772B75829 @ localhost";
Exec -ScriptBlock { gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B"; }
echo "Unlocked: F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B @ localhost";
Exec -ScriptBlock { gpg-preset-passphrase --passphrase "$unsecurePasswordProfessional" --preset "7F2D9FFF2E1D3A21299052552E7F68C82CD71C86"; }
echo "Unlocked: 7F2D9FFF2E1D3A21299052552E7F68C82CD71C86 @ localhost";
Exec -ScriptBlock { gpg-preset-passphrase --passphrase "$unsecurePasswordProfessional" --preset "5C31B095A9E5904D20A547DCF7E5096196D54909"; }
echo "Unlocked: 5C31B095A9E5904D20A547DCF7E5096196D54909 @ localhost";

# NOTE: The windows ssh-agent is implemented through the Windows registry and
# thus once the keys are added they do not need to be added again, even after reboot.

# The following will wait untill our dev-server has started
echo "Wait for dev-server.hyper-v.local";
RetryCommand -Verbose -ScriptBlock {
	ssh dev-server true;
	if ($LastExitCode -ne 0) {
		throw "failed";
	}
}

# Now unlock all the keys inside our dev-server
Exec -ScriptBlock { ssh dev-server unlock-ssh-key "~/.ssh/brad@bjc.id.au" "'$unsecurePasswordPersonal'"; }
Exec -ScriptBlock { ssh dev-server unlock-ssh-key "~/.ssh/brad.jones@xero.com" "'$unsecurePasswordProfessional'"; }
Exec -ScriptBlock { ssh dev-server unlock-gpg-key "83D182028C7F2DF102F09E61FF308BBB10F539D8" "'$unsecurePasswordPersonal'"; }
echo "Unlocked: 83D182028C7F2DF102F09E61FF308BBB10F539D8 @ dev-server.hyper-v.local";
Exec -ScriptBlock { ssh dev-server unlock-gpg-key "F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB" "'$unsecurePasswordPersonal'"; }
echo "Unlocked: F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB @ dev-server.hyper-v.local";
Exec -ScriptBlock { ssh dev-server unlock-gpg-key "1A8059A4CC0F06F670492ABBD0053F0772B75829" "'$unsecurePasswordPersonal'"; }
echo "Unlocked: 1A8059A4CC0F06F670492ABBD0053F0772B75829 @ dev-server.hyper-v.local";
Exec -ScriptBlock { ssh dev-server unlock-gpg-key "F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B" "'$unsecurePasswordPersonal'"; }
echo "Unlocked: F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B @ dev-server.hyper-v.local";
Exec -ScriptBlock { ssh dev-server unlock-gpg-key "7F2D9FFF2E1D3A21299052552E7F68C82CD71C86" "'$unsecurePasswordProfessional'"; }
echo "Unlocked: 7F2D9FFF2E1D3A21299052552E7F68C82CD71C86 @ dev-server.hyper-v.local";
Exec -ScriptBlock { ssh dev-server unlock-gpg-key "5C31B095A9E5904D20A547DCF7E5096196D54909" "'$unsecurePasswordProfessional'"; }
echo "Unlocked: 5C31B095A9E5904D20A547DCF7E5096196D54909 @ dev-server.hyper-v.local";