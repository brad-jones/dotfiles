$Stt = New-ScheduledTaskTrigger -AtLogOn;
		
$Sta = New-ScheduledTaskAction `
    -Execute "C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe" `
    -Argument "-NoLogo -NoProfile -File .\run-at-logon.ps1" `
    -WorkingDirectory "$env:USERPROFILE\Documents\WindowsPowershell\Scripts";

$u = whoami;
$STPrincipal = New-ScheduledTaskPrincipal `
    -UserID "$u" `
    -LogonType Interactive `
    -RunLevel Highest;

Register-ScheduledTask "Run at Logon" `
    -Principal $STPrincipal `
    -Trigger $Stt `
    -Action $Sta;