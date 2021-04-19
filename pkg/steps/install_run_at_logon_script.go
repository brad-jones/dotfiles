package steps

import (
	"fmt"
	"runtime"

	"github.com/brad-jones/dotfiles/pkg/tools/winsudo"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/brad-jones/gopwsh"
)

var ScheduledTaskError = goerr.New("failed")

func MustInstallRunAtLogonScript() {
	prefix := colorchooser.Sprint("install-run-at-logon-script")

	if runtime.GOOS == "windows" {
		ps := gopwsh.MustNew(gopwsh.Elevated(winsudo.Path()))
		defer ps.Exit()

		fmt.Println(prefix, "| does task exist?")
		if _, stderr := ps.MustExecute(`Get-ScheduledTask -TaskName "Run at Logon"`); len(stderr) == 0 {
			fmt.Println(prefix, "| deleting task...")
			if _, stderr := ps.MustExecute(`Unregister-ScheduledTask -TaskName "Run at Logon" -Confirm:$false`); len(stderr) > 0 {
				goerr.Check(ScheduledTaskError, "failed to unregister task", stderr)
			}
			fmt.Println(prefix, "| re-creating...")
		} else {
			fmt.Println(prefix, "| task does not exist, creating...")
		}

		if _, stderr := ps.MustExecute(
			`$u = whoami;`,
			`$Stt = New-ScheduledTaskTrigger -AtLogOn -User "$u";`,
			`$Sta = New-ScheduledTaskAction -Execute "C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe" -Argument "-NoLogo .\run-at-logon.ps1" -WorkingDirectory "$env:USERPROFILE\Documents\WindowsPowershell\Scripts";`,
			`$STPrincipal = New-ScheduledTaskPrincipal -UserID "$u" -LogonType Interactive;`,
			`Register-ScheduledTask "Run at Logon" -Principal $STPrincipal -Trigger $Stt -Action $Sta;`,
		); len(stderr) > 0 {
			goerr.Check(ScheduledTaskError, "failed to create task", stderr)
		}

		fmt.Println(prefix, "| task created")
	}
}

func InstallRunAtLogonScriptAsync() *task.Task {
	return task.New(func() { MustInstallRunAtLogonScript() })
}
