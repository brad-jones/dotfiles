. $env:USERPROFILE\Documents\WindowsPowerShell\utils.ps1;

# Make sure the windows GPG agent is running
Write-Output "Start the GPG agent";
RetryCommand -Verbose -ScriptBlock {
	gpg-connect-agent /bye;
	if ($LastExitCode -ne 0) {
		throw "failed";
	}
}

# Unlock personal keys
Write-Output "Getting passphrase:brad@bjc.id.au";
$passwordPersonal = Get-StoredCredential -Target "passphrase:brad@bjc.id.au";
$unsecurePasswordPersonal = [System.Net.NetworkCredential]::new('', $passwordPersonal.Password).Password;
gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "83D182028C7F2DF102F09E61FF308BBB10F539D8";
Write-Output "Unlocked: 83D182028C7F2DF102F09E61FF308BBB10F539D8 @ localhost";
gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB";
Write-Output "Unlocked: F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB @ localhost";
gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "1A8059A4CC0F06F670492ABBD0053F0772B75829";
Write-Output "Unlocked: 1A8059A4CC0F06F670492ABBD0053F0772B75829 @ localhost";
gpg-preset-passphrase --passphrase "$unsecurePasswordPersonal" --preset "F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B";
Write-Output "Unlocked: F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B @ localhost";

# NOTE: The windows ssh-agent is implemented through the Windows registry and
# thus once the keys are added they do not need to be added again, even after reboot.
# However some windows apps (looking at you GitKraken) do not integrate with
# the OpenSSH agent at all and prefer using Putty's pageant.exe
gopass bin cp "keys/ssh/brad@bjc.id.au.ppk" "$env:USERPROFILE\.ssh\brad@bjc.id.au.ppk";
pageant "$env:USERPROFILE\.ssh\brad@bjc.id.au.ppk";
rm -Force "$env:USERPROFILE\.ssh\brad@bjc.id.au.ppk";

# Unlock professional keys
if ($env:COMPUTERNAME -eq "XLW-5CD936CWNQ") {
	Write-Output "Getting passphrase:brad.jones@xero.com";
	$passwordProfessional = Get-StoredCredential -Target "passphrase:brad.jones@xero.com";
	$unsecurePasswordProfessional = [System.Net.NetworkCredential]::new('', $passwordProfessional.Password).Password;
	gpg-preset-passphrase --passphrase "$unsecurePasswordProfessional" --preset "7F2D9FFF2E1D3A21299052552E7F68C82CD71C86";
	Write-Output "Unlocked: 7F2D9FFF2E1D3A21299052552E7F68C82CD71C86 @ localhost";
	gpg-preset-passphrase --passphrase "$unsecurePasswordProfessional" --preset "5C31B095A9E5904D20A547DCF7E5096196D54909";
	Write-Output "Unlocked: 5C31B095A9E5904D20A547DCF7E5096196D54909 @ localhost";

	gopass bin cp "keys/ssh/brad.jones@xero.com.ppk" "$env:USERPROFILE\.ssh\brad.jones@xero.com.ppk";
	pageant "$env:USERPROFILE\.ssh\brad.jones@xero.com.ppk";
	rm -Force "$env:USERPROFILE\.ssh\brad.jones@xero.com.ppk";

	Write-Output "Wait for dev-server";
	RetryCommand -Verbose -ScriptBlock {
		ssh dev-server true;
		if ($LastExitCode -ne 0) {
			throw "failed";
		}
	}
	
	ssh dev-server unlock-ssh-key "~/.ssh/brad@bjc.id.au" "'$unsecurePasswordPersonal'";
	ssh dev-server unlock-ssh-key "~/.ssh/brad.jones@xero.com" "'$unsecurePasswordProfessional'";
	ssh dev-server unlock-gpg-key "83D182028C7F2DF102F09E61FF308BBB10F539D8" "'$unsecurePasswordPersonal'";
	Write-Output "Unlocked: 83D182028C7F2DF102F09E61FF308BBB10F539D8 @ dev-server";
	ssh dev-server unlock-gpg-key "F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB" "'$unsecurePasswordPersonal'";
	Write-Output "Unlocked: F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB @ dev-server";
	ssh dev-server unlock-gpg-key "1A8059A4CC0F06F670492ABBD0053F0772B75829" "'$unsecurePasswordPersonal'";
	Write-Output "Unlocked: 1A8059A4CC0F06F670492ABBD0053F0772B75829 @ dev-server";
	ssh dev-server unlock-gpg-key "F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B" "'$unsecurePasswordPersonal'";
	Write-Output "Unlocked: F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B @ dev-server";
	ssh dev-server unlock-gpg-key "7F2D9FFF2E1D3A21299052552E7F68C82CD71C86" "'$unsecurePasswordProfessional'";
	Write-Output "Unlocked: 7F2D9FFF2E1D3A21299052552E7F68C82CD71C86 @ dev-server";
	ssh dev-server unlock-gpg-key "5C31B095A9E5904D20A547DCF7E5096196D54909" "'$unsecurePasswordProfessional'";
	Write-Output "Unlocked: 5C31B095A9E5904D20A547DCF7E5096196D54909 @ dev-server";
}