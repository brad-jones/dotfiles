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
RetryCommand -Verbose -ScriptBlock {
	gpg-connect-agent /bye;
	if ($LastExitCode -ne 0) {
		throw "failed";
	}
}

# Get our passphrases from the windows cred store
$passwordPersonal = Get-StoredCredential -Target "passphrase:brad@bjc.id.au";
$unsecurePasswordPersonal = [System.Net.NetworkCredential]::new('', $passwordPersonal.Password).Password;
$passwordProfessional = Get-StoredCredential -Target "passphrase:brad.jones@xero.com";
$unsecurePasswordProfessional = [System.Net.NetworkCredential]::new('', $passwordProfessional.Password).Password;

# Add the GPG keys to the windows agent
gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "83D182028C7F2DF102F09E61FF308BBB10F539D8";
gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB";
gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "1A8059A4CC0F06F670492ABBD0053F0772B75829";
gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B";
gpg-preset-passphrase --passphrase "$unsecurePasswordProfessional" --preset "7F2D9FFF2E1D3A21299052552E7F68C82CD71C86";
gpg-preset-passphrase --passphrase "$unsecurePasswordProfessional" --preset "5C31B095A9E5904D20A547DCF7E5096196D54909";

# NOTE: The windows ssh-agent is implemented through the Windows registry and
# thus once the keys are added they do not need to be added again, even after reboot.

# The following will wait untill our dev-server has started
Retry-Command -Verbose -ScriptBlock {
	ssh dev-server true;
	if ($LastExitCode -ne 0) {
		throw "failed";
	}
}

# Now unlock all the keys inside our dev-server
ssh dev-server unlock-ssh-key "~/.ssh/brad@bjc.id.au" "'$unsecurePasswordPersonal'";
ssh dev-server unlock-ssh-key "~/.ssh/brad.jones@xero.com" "'$unsecurePasswordProfessional'";
ssh dev-server unlock-gpg-key "83D182028C7F2DF102F09E61FF308BBB10F539D8" "'$unsecurePasswordPersonal'";
ssh dev-server unlock-gpg-key "F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB" "'$unsecurePasswordPersonal'";
ssh dev-server unlock-gpg-key "1A8059A4CC0F06F670492ABBD0053F0772B75829" "'$unsecurePasswordPersonal'";
ssh dev-server unlock-gpg-key "F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B" "'$unsecurePasswordPersonal'";
ssh dev-server unlock-gpg-key "7F2D9FFF2E1D3A21299052552E7F68C82CD71C86" "'$unsecurePasswordProfessional'";
ssh dev-server unlock-gpg-key "5C31B095A9E5904D20A547DCF7E5096196D54909" "'$unsecurePasswordProfessional'";