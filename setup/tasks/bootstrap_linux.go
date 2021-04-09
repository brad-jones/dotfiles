package tasks

import (
	"github.com/brad-jones/dotfiles/setup/tasks/steps"
	"github.com/brad-jones/goasync/v2/await"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
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
		steps.InstallGitGpgAsync(answers.SudoPassword),
		steps.InstallGithubPkgAsync("gopasspw", "gopass", "v1.9.2", "", "gopass", ""),
		steps.InstallGithubPkgAsync("twpayne", "chezmoi", "v1.8.8", "", "chezmoi", ""),
	)

	await.MustFastAllOrError(
		steps.ChezmoiInitAsync(answers.GithubPassword),
		steps.DownloadVaultAsync(answers.GithubPassword),
		steps.DownloadVaultKeyAsync(answers.GitlabPassword, answers.VaultKeyPassword),
	)

	goexec.MustRun("chezmoi", "apply", "--debug")

	return
}
