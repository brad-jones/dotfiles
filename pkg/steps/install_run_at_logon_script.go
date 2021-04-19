package steps

import (
	"fmt"
	"os"
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

	exe, err := os.Executable()
	goerr.Check(err, "failed to get path to current running exe")

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
			`$Sta = New-ScheduledTaskAction -Execute "`+exe+`" -Argument "-update" -WorkingDirectory "$env:USERPROFILE";`,
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
