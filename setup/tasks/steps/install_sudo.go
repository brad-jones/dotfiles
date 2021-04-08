// +build windows

package steps

import (
	"github.com/brad-jones/goasync/v2/task"
)

// MustInstallSudoForWindows installs sudo for Windows.
//
// Subsequent steps may need to elevate. On *nix this is of course built-in.
// On Windows we need to install a sudo like tool. We can't really use scoop
// to install the tool because scoop it's self needs to elevate on install.
//
// There are 2 options:
// - https://github.com/brad-jones/winsudo
//   My tool that I built first in golang.
//   It does work but has some edge case bugs that need to be resolved.
//
// - https://github.com/gerardog/gsudo
//   This is a .NET Framework tool that I found later and is probably more
//   robust, better tested, more mature & smarter with respect to how it
//   actually does the elevation.
func MustInstallSudoForWindows() {
	//MustInstallGithubPkg("brad-jones", "winsudo", "v1.0.5", `winsudo_amd64\.zip`, "sudo", "")
	MustInstallGithubPkg("gerardog", "gsudo", "v0.7.3", `gsudo\..*\.zip`, "gsudo", "sudo")
}

func InstallSudoForWindowsAsync() *task.Task {
	return task.New(func() { MustInstallSudoForWindows() })
}
