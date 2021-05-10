package gopass

import (
	"encoding/base64"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

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
		goexec.MustRunBuffered(Path(), "show", "-o", key).StdOut,
	)
}

// Returns a binary secret from the vault
func GetBinarySecret(key string) []byte {
	out := goexec.MustRunBuffered(Path(), "show", key+".b64").StdOut
	b64 := strings.TrimSpace(strings.Join(strings.Split(out, "\n")[1:], ""))
	plainTxt, err := base64.StdEncoding.DecodeString(b64)
	goerr.Check(err, "failed to decode base64 bytes")
	return plainTxt
}

func WriteBinarySecret(key, dst string, perm fs.FileMode) {
	goerr.Check(
		ioutil.WriteFile(dst, GetBinarySecret(key), perm),
		"failed to write secret", key, dst,
	)
}

// Indicates if the vault is working & unlocked or not
func VaultUnlocked() (r bool) {
	defer goerr.Handle(func(err error) { r = false })
	if !utils.CommandExists(Path()) {
		return false
	}
	return len(GetSecret("ato/tfn")) > 0
}
