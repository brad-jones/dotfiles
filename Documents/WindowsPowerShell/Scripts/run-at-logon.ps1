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

# Make sure the windows GPG agent is running
echo "Start the GPG agent";
RetryCommand -Verbose -ScriptBlock {
	gpg-connect-agent /bye;
	if ($LastExitCode -ne 0) {
		throw "failed";
	}
}

# Unlock personal keys
echo "Getting passphrase:brad@bjc.id.au";
$passwordPersonal = Get-StoredCredential -Target "passphrase:brad@bjc.id.au";
$unsecurePasswordPersonal = [System.Net.NetworkCredential]::new('', $passwordPersonal.Password).Password;
gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "83D182028C7F2DF102F09E61FF308BBB10F539D8";
echo "Unlocked: 83D182028C7F2DF102F09E61FF308BBB10F539D8 @ localhost";
gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB";
echo "Unlocked: F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB @ localhost";
gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "1A8059A4CC0F06F670492ABBD0053F0772B75829";
echo "Unlocked: 1A8059A4CC0F06F670492ABBD0053F0772B75829 @ localhost";
gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B";
echo "Unlocked: F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B @ localhost";

# NOTE: The windows ssh-agent is implemented through the Windows registry and
# thus once the keys are added they do not need to be added again, even after reboot.

# Unlock professional keys
if ($env:COMPUTERNAME -eq "XLW-5CD936CWNQ") {
	echo "Getting passphrase:brad.jones@xero.com";
	$passwordProfessional = Get-StoredCredential -Target "passphrase:brad.jones@xero.com";
	$unsecurePasswordProfessional = [System.Net.NetworkCredential]::new('', $passwordProfessional.Password).Password;
	gpg-preset-passphrase --passphrase "$unsecurePasswordProfessional" --preset "7F2D9FFF2E1D3A21299052552E7F68C82CD71C86";
	echo "Unlocked: 7F2D9FFF2E1D3A21299052552E7F68C82CD71C86 @ localhost";
	gpg-preset-passphrase --passphrase "$unsecurePasswordProfessional" --preset "5C31B095A9E5904D20A547DCF7E5096196D54909";
	echo "Unlocked: 5C31B095A9E5904D20A547DCF7E5096196D54909 @ localhost";
	echo "Wait for dev-server-hv.local";
	RetryCommand -Verbose -ScriptBlock {
		ssh dev-server-hv true;
		if ($LastExitCode -ne 0) {
			throw "failed";
		}
	}
	ssh dev-server-hv unlock-ssh-key "~/.ssh/brad@bjc.id.au" "'$unsecurePasswordPersonal'";
	ssh dev-server-hv unlock-ssh-key "~/.ssh/brad.jones@xero.com" "'$unsecurePasswordProfessional'";
	ssh dev-server-hv unlock-gpg-key "83D182028C7F2DF102F09E61FF308BBB10F539D8" "'$unsecurePasswordPersonal'";
	echo "Unlocked: 83D182028C7F2DF102F09E61FF308BBB10F539D8 @ dev-server-hv.local";
	ssh dev-server-hv unlock-gpg-key "F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB" "'$unsecurePasswordPersonal'";
	echo "Unlocked: F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB @ dev-server-hv.local";
	ssh dev-server-hv unlock-gpg-key "1A8059A4CC0F06F670492ABBD0053F0772B75829" "'$unsecurePasswordPersonal'";
	echo "Unlocked: 1A8059A4CC0F06F670492ABBD0053F0772B75829 @ dev-server-hv.local";
	ssh dev-server-hv unlock-gpg-key "F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B" "'$unsecurePasswordPersonal'";
	echo "Unlocked: F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B @ dev-server-hv.local";
	ssh dev-server-hv unlock-gpg-key "7F2D9FFF2E1D3A21299052552E7F68C82CD71C86" "'$unsecurePasswordProfessional'";
	echo "Unlocked: 7F2D9FFF2E1D3A21299052552E7F68C82CD71C86 @ dev-server-hv.local";
	ssh dev-server-hv unlock-gpg-key "5C31B095A9E5904D20A547DCF7E5096196D54909" "'$unsecurePasswordProfessional'";
	echo "Unlocked: 5C31B095A9E5904D20A547DCF7E5096196D54909 @ dev-server-hv.local";
}