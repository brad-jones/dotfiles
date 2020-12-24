package tasks

import (
	"runtime"

	"github.com/brad-jones/dotfiles/setup/tasks/steps"
	"github.com/brad-jones/goasync/v2/await"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

// Bootstrap will run when someone executes this program directly without
// any additional input and it's job is to do all the things required to
// setup chezmoi.
//
// https://github.com/brad-jones/dotfiles/blob/68aaced5cc2e67007e7c0024adc0142cdd95502e/install.sh
// https://github.com/brad-jones/dotfiles/blob/68aaced5cc2e67007e7c0024adc0142cdd95502e/install.ps1
func Bootstrap() (err error) {
	defer goerr.Handle(func(e error) { err = e })

	answers := steps.BootstrapSurvey()

	await.MustFastAllOrError(
		task.New(func() { steps.InstallGitGpg(answers.SudoPassword) }),
		task.New(func() { steps.InstallGithubPkg("gopasspw", "gopass", "v1.9.2") }),
		task.New(func() { steps.InstallGithubPkg("twpayne", "chezmoi", "v1.8.8") }),
	)

	await.MustFastAllOrError(
		task.New(func() { steps.ChezmoiInit(answers.GithubPassword) }),
		task.New(func() { steps.DownloadVault(answers.GithubPassword) }),
		task.New(func() { steps.DownloadVaultKey(answers.GitlabPassword, answers.VaultKeyPassword) }),
	)

	if runtime.GOOS == "windows" {
		prefix := colorchooser.Sprint("cleanup-wincred")
		await.MustFastAllOrError(
			task.New(func() { goexec.MustRunPrefixed(prefix, "cmdkey", `/delete:git:https://github.com`) }),
			task.New(func() { goexec.MustRunPrefixed(prefix, "cmdkey", `/delete:git:https://brad-jones@gitlab.com`) }),
		)
	}

	goexec.MustRun("chezmoi", "apply", "--debug")

	return
}
