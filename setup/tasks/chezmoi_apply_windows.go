package tasks

import (
	"path/filepath"

	"github.com/brad-jones/dotfiles/setup/tasks/steps"
	"github.com/brad-jones/dotfiles/setup/tasks/utils"
	"github.com/brad-jones/goasync/v2/await"
	"github.com/brad-jones/goerr/v2"
)

// ChezmoiApply will run when chezmoi it's self executes
// this program via the `run_setup` scripts.
func ChezmoiApply() (err error) {
	defer goerr.Handle(func(e error) { err = e })

	// Make sure we have a way to elevate on Windows
	steps.MustInstallSudoForWindows()
	steps.MustDisableUACForWindows()

	// On Windows we install/update the majority of our tools via https://scoop.sh/
	steps.MustInstallScoop()
	await.MustFastAllOrError(
		steps.InstallScoopBucketAsync("extras", ""),
		steps.InstallScoopBucketAsync("nonportable", ""),
		steps.InstallScoopBucketAsync("java", ""),
		steps.InstallScoopBucketAsync("jetbrains", ""),
		steps.InstallScoopBucketAsync("goreleaser", "https://github.com/goreleaser/scoop-bucket.git"),
		steps.InstallScoopBucketAsync("brad-jones", "https://github.com/brad-jones/scoop-bucket.git"),
	)

	// Scoop makes use of aria2 to speed up downloads, so may as well have
	// it installed before installing the rest of the tools.
	// see: https://github.com/lukesampson/scoop#multi-connection-downloads-with-aria2
	steps.MustInstallScoopPkg("aria2", "")

	// Kill running apps the might get updated
	utils.KillProcByName("gpg-agent")
	utils.KillProcByName("ssh-agent")
	utils.KillProcByName("pageant")

	// [pkg]ver - if ver == "", then latest version will install
	steps.MustInstallScoopPkgs(map[string]string{
		"7zip":                 "",
		"adoptopenjdk-hotspot": "",
		"aws-vault":            "5.4.4",
		"aws":                  "",
		"curl":                 "",
		"dart":                 "",
		"deno":                 "",
		"git":                  "",
		"gitkraken":            "",
		"go":                   "",
		"gpg":                  "",
		"grep":                 "",
		"jq":                   "",
		"kotlin":               "",
		"ktlint":               "",
		"maven":                "",
		"nodejs":               "",
		"nuget":                "",
		"openssl":              "",
		"packer":               "",
		"protobuf":             "",
		"putty":                "",
		"pwsh":                 "",
		"python":               "",
		"ruby":                 "",
		"sed":                  "",
		"sonar-scanner":        "",
		"task":                 "",
		"terraform":            "",
		"vlc":                  "",
		"vscode":               "",
		"wget":                 "",
		"win32-openssh":        "",
		"windows-terminal":     "",
	})

	// Install the RunAtLogon script that auto unlocks our gopass vault &
	// SSH/GPG keys by piping passphrases stored in the WinCred store.
	steps.MustInstallCredentialManager()
	steps.MustInstallRunAtLoginPwshScript()

	// Chezmoi has this annoying habbit of setting the read-only bit which
	// makes Windows Terminal throw a warning on startup
	utils.SetWritable(filepath.Join(
		utils.HomeDir(),
		"AppData",
		"Local",
		"Microsoft",
		"Windows Terminal",
		"settings.json",
	))

	// Here we configure our primary WSL instance. Essentially this will
	// install a WSL distro and then run this same boostrap setup inside it.
	steps.MustInstallWSL(false)
	distro := steps.MustInstallWSLFedora("33.20210226.0", true, false)
	steps.MustBoostrapInWSL(distro)

	// These tasks are cross platform
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
