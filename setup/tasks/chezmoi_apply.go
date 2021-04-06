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

	steps.MustInstallSudoForWindows()
	steps.MustDisableUACForWindows()

	// Now install or update our SSH/GPG keys
	await.MustFastAllOrError(
		steps.InstallSSHGpgKeysAsync(),
		steps.InstallGithubPkgAsync("brad-jones", "ssh-add-with-pass", "v1.0.4", "", "ssh_add_with_pass", ""),
	)

	if runtime.GOOS == "windows" {
		// On Windows we install/update the majority of our tools via scoop
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
			"vscode":               "",
			"wget":                 "",
			"win32-openssh":        "",
			"windows-terminal":     "",
		})

		// Install the RunAtLogon script that auto unlocks our gopass vault &
		// SSH/GPG keys by piping passphrases stored in the WinCred store.
		steps.MustInstallCredentialManager()
		steps.MustInstallRunAtLoginPwshScript()

		// Highly likely to fail due to all this shenanigans:
		// - https://github.com/PowerShell/PSReadLine#upgrading
		// - https://github.com/PowerShell/PSReadLine/issues/1370
		//
		// Probs just going to leave this disabled for now as there
		// is no easy way to solve this.
		//steps.InstallPSReadLine()

		// We have these scripts written in dart under: ".local/sbin".
		//
		// They are tools which are very specific to me, do questionable things
		// like (automatically filling out MFA tokens) or are very much prototype
		// quality so they make less sense to release publically.
		//
		// The only issue is that when Dart updates then the scripts start to
		// break. I could "pin" the version of dart in this very setup project
		// but then I couldn't use the latest Dart on some other project, at
		// least on Windows anyway, where a dartenv style version manager isn't
		// really viable.
		//
		// The whole pattern is kinda cool, but kinda a pain at the same time.
		// I think longer term these things will turn into compiled self-contained
		// executables, written with either Dart, Go or even Deno perhaps...
		steps.MustInstallDartScriptDeps()
	}

	return
}

/*
	TODO

	Google Chrome, Firefox, Wavebox & other apps we want to install from the
	nonportable bucket & then basically just rely on the built-in updating.

	Fix Dart warnings

	Fix aws-vault rotate

	The rdp shortcuts

	Dotnet Core install, need a dotnet version manager.
	I think the install scripts actually do an ok job

	Docker Desktop

	WSL setup, which then leads into Linux setup

	Vscode extensions sync
*/
