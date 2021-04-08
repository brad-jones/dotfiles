// +build windows

package steps

import (
	"fmt"

	"github.com/brad-jones/dotfiles/setup/tasks/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/brad-jones/gopwsh"
)

// MustInstallScoop will install the https://scoop.sh/ package manager for Windows
//
// TODO: change installation dir to ~/.scoop
func MustInstallScoop() {
	prefix := colorchooser.Sprint("install-scoop")

	ps := gopwsh.MustNew()
	defer ps.Exit()

	if _, stderr := ps.MustExecute("Get-Command scoop"); len(stderr) == 0 {
		goexec.MustRunPrefixed(prefix, "powershell",
			"-Command", "scoop update",
		)
		return
	}

	psElevated := gopwsh.MustNew(gopwsh.Elevated(utils.SudoBin()))
	defer psElevated.Exit()

	fmt.Println(prefix, "| setting execution policy to RemoteSigned")
	psElevated.MustExecute("Set-ExecutionPolicy RemoteSigned -scope CurrentUser")

	goexec.MustRunPrefixed(prefix, "powershell",
		"-Command", "iwr -useb get.scoop.sh | iex",
	)
}

func InstallScoopAsync() *task.Task {
	return task.New(func() { MustInstallScoop() })
}
