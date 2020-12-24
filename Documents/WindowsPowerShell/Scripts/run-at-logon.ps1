. $env:USERPROFILE\Documents\WindowsPowerShell\utils.ps1;

# Start the ssh agent
Write-Output "Start the SSH agent";
Start-Service ssh-agent;

# Start the gpg agent
Write-Output "Start the GPG agent";
RetryCommand -Verbose -ScriptBlock {
	gpg-connect-agent /bye;
	if ($LastExitCode -ne 0) {
		throw "failed";
	}
}

# Unlock gopass vault on login
Write-Output "Unlock the gopass vault";
$vaultPass = Get-StoredCredential -Target "passphrase:vault";
$unsecureVaultPass = [System.Net.NetworkCredential]::new('', $vaultPass.Password).Password;
gpg-preset-passphrase --passphrase "$unsecureVaultPass" --preset "83D182028C7F2DF102F09E61FF308BBB10F539D8";
gpg-preset-passphrase --passphrase "$unsecureVaultPass" --preset "F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB";
Write-Output "Unlocked gpg key: Brad Jones (vault) <brad@bjc.id.au>";

# Unlock personal gpg key
$pass = gopass show "keys/gpg/brad@bjc.id.au.pass";
gpg-preset-passphrase --passphrase "$pass" --preset "1A8059A4CC0F06F670492ABBD0053F0772B75829";
gpg-preset-passphrase --passphrase "$pass" --preset "F1C1E6443BB1B7AA8062DF0E085C64B391E94D5B";
Write-Output "Unlocked gpg key: Brad Jones <brad@bjc.id.au>";

# Unlock personal ssh key
$pass = gopass show "keys/ssh/brad@bjc.id.au.pass";
Write-Output $pass | ssh-add-with-pass "$env:USERPROFILE\.ssh\brad@bjc.id.au";
Write-Output "Unlocked ssh key: brad@bjc.id.au";

# Some windows apps (looking at you GitKraken) do not integrate with
# the OpenSSH agent at all and prefer using Putty's pageant.exe
gopass bin cp "keys/ssh/brad@bjc.id.au.ppk" "$env:USERPROFILE\.ssh\brad@bjc.id.au.ppk";
pageant "$env:USERPROFILE\.ssh\brad@bjc.id.au.ppk"; Start-Sleep -s 1;
Remove-Item -Force "$env:USERPROFILE\.ssh\brad@bjc.id.au.ppk";
Write-Output "Added brad@bjc.id.au to pageant";

# Unlock professional keys
if ($env:COMPUTERNAME -eq "XLW-5CD936CWNQ") {
	$pass = gopass show "keys/gpg/brad.jones@xero.com.pass";
	gpg-preset-passphrase --passphrase "$pass" --preset "7F2D9FFF2E1D3A21299052552E7F68C82CD71C86";
	gpg-preset-passphrase --passphrase "$pass" --preset "5C31B095A9E5904D20A547DCF7E5096196D54909";
	Write-Output "Unlocked gpg key: Brad Jones <brad.jones@xero.com>";

	$pass = gopass show "keys/ssh/brad.jones@xero.com.pass";
	Write-Output $pass | ssh-add-with-pass "$env:USERPROFILE\.ssh\brad.jones@xero.com.pass";
	Write-Output "Unlocked ssh key: brad.jones@xero.com";

	gopass bin cp "keys/ssh/brad.jones@xero.com.ppk" "$env:USERPROFILE\.ssh\brad.jones@xero.com.ppk";
	pageant "$env:USERPROFILE\.ssh\brad.jones@xero.com.ppk"; Start-Sleep -s 1;
	Remove-Item -Force "$env:USERPROFILE\.ssh\brad.jones@xero.com.ppk";
	Write-Output "Added brad.jones@xero.com to pageant";
}