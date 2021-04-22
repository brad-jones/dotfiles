package winsudo

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/brad-jones/dotfiles/pkg/tools"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/dotfiles/pkg/utils/ghpkg"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/brad-jones/gopwsh"
)

// NotWindows is returned whenever one of these tasks are run on a non-windows OS
var NotWindows = goerr.New("winsudo is a windows only technology")

// Installs sudo for Windows.
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
//   actually does the elevation. So running with this for now...
//
// NOTE: My golang tool is faster :)
func Install(reset bool) (err error) {
	defer goerr.Handle(func(e error) { err = e })
	if runtime.GOOS != "windows" {
		goerr.Check(NotWindows)
	}

	exists := utils.CommandExists("sudo")

	v := tools.GetVersion("gsudo")
	goerr.Check(ghpkg.InstallPkg("gerardog", "gsudo", v.No,
		ghpkg.ExeName("gsudo"),
		ghpkg.DstExeName("sudo"),
		ghpkg.PkgPattern(`gsudo\..*\.zip`),
		ghpkg.Sha256Hash(v.Hash),
		ghpkg.Reset(reset),
	), "failed to install gsudo")

	// I liken UAC to SELinux, I have gotten by without both technologies and
	// I am yet to have a security incident. Perhaps thats naive but oh well...
	//
	// This does not disable UAC completely, we do not set `EnableLUA` to 0.
	// As this would break vital built-in components of Windows,
	// like the Firewall.
	//
	// We just stop the prompt from displaying for Admin users with
	// `ConsentPromptBehaviorAdmin` set to 0.
	//
	// Actually a more accurate analogy is where Sudo prompts for a password.
	// This is essentially like having `ALL=(ALL) NOPASSWD:ALL`. If UAC had some
	// sort of CLI interface this wouldn't be needed.
	if reset || !exists {
		prefix := colorchooser.Sprint("install-gsudo")
		fmt.Println(prefix, "|", "setting ConsentPromptBehaviorAdmin to 0")
		ps := gopwsh.MustNew(gopwsh.Elevated(Path()))
		defer ps.Exit()
		ps.MustExecute(`Set-ItemProperty -Path REGISTRY::HKEY_LOCAL_MACHINE\Software\Microsoft\Windows\CurrentVersion\Policies\System -Name ConsentPromptBehaviorAdmin -Value 0`)
		fmt.Println(prefix, "|", "ConsentPromptBehaviorAdmin is now 0")
	}

	return
}

// MustInstall does the same thing as Install but panics instead of returning an error
func MustInstall(reset bool) {
	goerr.Check(Install(reset))
}

// InstallAsync does the same thing as Install but asynchronously.
func InstallAsync(reset bool) *task.Task {
	return task.New(func() { MustInstall(reset) })
}

// Returns the path to the sudo binary
func Path() string {
	return filepath.Join(utils.HomeDir(), ".local", "bin", "sudo.exe")
}
