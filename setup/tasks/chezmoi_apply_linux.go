package tasks

import (
	"github.com/brad-jones/dotfiles/setup/tasks/steps"
	"github.com/brad-jones/goasync/v2/await"
	"github.com/brad-jones/goerr/v2"
)

// ChezmoiApply will run when chezmoi it's self executes
// this program via the `run_setup` scripts.
func ChezmoiApply() (err error) {
	defer goerr.Handle(func(e error) { err = e })

	await.MustFastAllOrError(
		steps.InstallSSHGpgKeysAsync(),
		steps.InstallChromeAsync(),
		steps.InstallFirefoxAsync(),
		steps.InstallWaveboxAsync(),
		steps.InstallDartScriptDepsAsync(),
		steps.InstallDotnetAsync("latest", "3.1.407", "2.1.814"),
		steps.InstallGithubPkgAsync("brad-jones", "ssh-add-with-pass", "v1.0.4", "", "ssh_add_with_pass", ""),
	)

	return
}
