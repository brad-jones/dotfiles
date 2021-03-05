package tasks

import (
	"runtime"

	"github.com/brad-jones/dotfiles/setup/tasks/steps"
	"github.com/brad-jones/goerr/v2"
)

// ChezmoiApply will run when chezmoi it's self executes
// this program via the `run_setup` scripts.
func ChezmoiApply() (err error) {
	defer goerr.Handle(func(e error) { err = e })

	steps.InstallSSHGpgKeys()
	steps.InstallGithubPkg("brad-jones", "ssh-add-with-pass", "v1.0.1", "ssh_add_with_pass")

	if runtime.GOOS == "windows" {
		steps.InstallScoop()
		steps.InstallScoopPkg("pwsh", "7.1.2")
		steps.InstallCredentialManager()
		//steps.InstallPSReadLine()
		//steps.InstallRunAtLoginPwshScript()
	}

	return nil
}
