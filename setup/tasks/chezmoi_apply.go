package tasks

import (
	"runtime"

	"github.com/brad-jones/dotfiles/setup/tasks/steps"
	"github.com/brad-jones/goasync/v2/await"
	"github.com/brad-jones/goerr/v2"
)

// ChezmoiApply will run when chezmoi it's self executes
// this program via the `run_setup` scripts.
func ChezmoiApply() (err error) {
	defer goerr.Handle(func(e error) { err = e })

	if runtime.GOOS == "windows" {
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
		steps.KillProcByName("gpg-agent")
		steps.KillProcByName("ssh-agent")
		steps.KillProcByName("pageant")

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
	}

	// These tasks are cross platform
	// We are probs being a bit ambitious trying to run all this is once... we will see
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

/*
	TODO

	install xop

	Fix Dart warnings

	Fix aws-vault rotate

	The rdp shortcuts

	Docker Desktop

	WSL setup, which then leads into Linux setup

	Vscode extensions sync
*/
