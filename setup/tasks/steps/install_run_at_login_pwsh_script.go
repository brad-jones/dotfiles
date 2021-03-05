package steps

// InstallRunAtLoginPwshScript makes our "run-at-logon.ps1" script run at login.
func InstallRunAtLoginPwshScript() {
	/*//prefix := colorchooser.Sprint("install-run-at-login-script")

	if !isElevated() {
		//Start-Process powershell -Wait -Verb RunAs `
		//	-ArgumentList "-NoLogo", "-NoProfile", `
		//	"-EncodedCommand", "${base64.encode(encodeUtf16le(script))}";
	}

	ps, err := powershell.New(&backend.Local{})
	goerr.Check(err)
	defer ps.Exit()

	_, _, err = ps.Execute("$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent());")
	goerr.Check(err)
	stdout, _, err := ps.Execute("$currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator);")
	goerr.Check(err)

	// \$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent());
	// \$currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator);

	fmt.Println("|" + stdout + "|")
	fmt.Println(stderr)
	//_, _, err = ps.Execute("Unregister-ScheduledTask -TaskName \"Run at Logon\" -Confirm:$false")
	//goerr.Check(err)*/
}

/*
# Install our run at login script
# ------------------------------------------------------------------------------
if (-Not (Get-ScheduledTask -TaskName "Run at Logon" -ErrorAction Ignore)) {
	Exec -ScriptBlock {
		$u = whoami;

		$Stt = New-ScheduledTaskTrigger -AtLogOn -User "$u";

		$Sta = New-ScheduledTaskAction `
			-Execute "C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe" `
			-Argument "-NoLogo .\run-at-logon.ps1" `
			-WorkingDirectory "$env:USERPROFILE\Documents\WindowsPowershell\Scripts";

		$STPrincipal = New-ScheduledTaskPrincipal `
			-UserID "$u" `
			-LogonType Interactive;

		Register-ScheduledTask "Run at Logon" `
			-Principal $STPrincipal `
			-Trigger $Stt `
			-Action $Sta;
	}
}
*/
