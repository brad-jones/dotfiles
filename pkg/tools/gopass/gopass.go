package gopass

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/brad-jones/dotfiles/pkg/assets"
	"github.com/brad-jones/dotfiles/pkg/tools"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/dotfiles/pkg/utils/ghpkg"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
)

// Installs the https://github.com/gopasspw/gopass tool
func Install(reset bool) (err error) {
	defer goerr.Handle(func(e error) { err = e })
	v := tools.GetVersion("gopass")
	goerr.Check(ghpkg.InstallPkg("gopasspw", "gopass", v.No,
		ghpkg.Sha256Hash(v.Hash),
		ghpkg.Reset(reset),
	), "failed to install gopass")
	assets.WriteFolderToHome("AppData/Local/gopass")
	assets.WriteFolderToHome(".config/gopass")
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

// Returns the path to the gopass binary
func Path() string {
	p := filepath.Join(utils.HomeDir(), ".local", "bin", "gopass")
	if runtime.GOOS == "windows" {
		p = p + ".exe"
	}
	return p
}

// Returns a secret from the vault
func GetSecret(key string) string {
	return strings.TrimSpace(
		goexec.MustRunBuffered("gopass", "show", "-o", key).StdOut,
	)
}

// Indicates if the vault is working & unlocked or not
func VaultUnlocked() (r bool) {
	defer goerr.Handle(func(err error) { r = false })
	if !utils.CommandExists("gopass") {
		return false
	}
	return len(GetSecret("ato/tfn")) > 0
}
