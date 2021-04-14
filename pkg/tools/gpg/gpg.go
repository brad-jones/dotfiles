package gpg

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/avast/retry-go"
	"github.com/brad-jones/dotfiles/pkg/assets"
	"github.com/brad-jones/dotfiles/pkg/tools/scoop"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/gosimple/slug"
)

// UnSupportedOS is returned when we tried to install on an OS we don't recognize
var UnSupportedOS = goerr.New("we do not support the OS you are using")

// Installs the gpg tool
func Install(sudoPassword string) (err error) {
	defer goerr.Handle(func(e error) { err = e })
	prefix := colorchooser.Sprint("install-gpg")

	if runtime.GOOS == "windows" {
		goerr.Check(scoop.InstallOrUpdatePkgs(map[string]string{"gpg": "*"}), "failed to install gpg")
		configure()
		return
	}

	if utils.CommandExists("dnf") {
		if utils.IsRoot() {
			goexec.MustRunPrefixed(prefix, "dnf", "install", "-y", "gnupg", "pinentry")
		} else {
			if len(sudoPassword) > 0 {
				goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd("sudo",
					goexec.SetIn(strings.NewReader(sudoPassword)),
					goexec.Args("-S", "dnf", "install", "-y", "gnupg", "pinentry"),
				))
			} else {
				goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd("sudo",
					goexec.Args("dnf", "install", "-y", "gnupg", "pinentry"),
				))
			}
		}
		configure()
		return
	}

	goerr.Check(UnSupportedOS)
	return
}

func configure() {
	if runtime.GOOS == "windows" {
		assets.WriteFolder(".gnupg",
			filepath.Join(scoop.Path(), "apps/gpg/current/home"),
		)
		return
	}

	assets.WriteFolderToHome(".gnupg")
}

// MustInstall does the same thing as Install but panics instead of returning an error
func MustInstall(sudoPassword string) {
	goerr.Check(Install(sudoPassword))
}

// InstallAsync does the same thing as Install but asynchronously.
func InstallAsync(sudoPassword string) *task.Task {
	return task.New(func() { MustInstall(sudoPassword) })
}

func MustStartAgent() {
	prefix := colorchooser.Sprint("gpg-agent")

	if runtime.GOOS == "windows" {
		fmt.Println(prefix, "|", "starting...")
		goerr.Check(retry.Do(func() error {
			return goexec.RunPrefixed(prefix, "gpg-connect-agent", "/bye")
		}))
		fmt.Println(prefix, "|", "running")
	}
}

func PublicKeyExists(keyName string) bool {
	_, err := goexec.RunBuffered("gpg", "-k", keyName)
	return err == nil
}

func PrivateKeyExists(keyName string) bool {
	_, err := goexec.RunBuffered("gpg", "-K", keyName)
	return err == nil
}

func MustImportKey(keyPath, keyName, keyPassphrase string, unlock bool) {
	prefix := colorchooser.Sprint("import-gpg-key-" + slug.Make(keyName))

	if !PrivateKeyExists(keyName) {
		fmt.Println(prefix, "|", "importing", keyPath)
		if len(keyPassphrase) > 0 {
			goexec.MustRunPrefixed(prefix, "gpg",
				"--pinentry-mode", "loopback",
				"--passphrase", keyPassphrase,
				"--import", keyPath,
			)
		} else {
			goexec.MustRunPrefixed(prefix, "gpg",
				"--import", keyPath,
			)
		}

		fmt.Println(prefix, "|", "imported", keyPath, "trusting", keyName)
		MustTrustKey(keyName)
		fmt.Println(prefix, "|", "trusted", keyName)
	}

	if unlock {
		MustUnlockKey(keyName, keyPassphrase)
	}
}

func MustTrustKey(keyName string) {
	goexec.MustRunBufferedCmd(goexec.MustCmd("gpg",
		goexec.SetIn(strings.NewReader("5\ny\n")),
		goexec.Args(
			"--command-fd", "0",
			"--edit-key", keyName,
			"trust",
		),
	))
}

func MustUnlockKey(keyName, keyPassphrase string) {
	prefix := colorchooser.Sprint("unlock-gpg-key-" + slug.Make(keyName))

	presetBin := "gpg-preset-passphrase"
	if runtime.GOOS == "linux" {
		presetBin = "/usr/libexec/gpg-preset-passphrase"
	}

	for _, line := range strings.Split(goexec.MustRunBuffered("gpg", "-K", "--with-keygrip", keyName).StdOut, "\n") {
		if !strings.Contains(line, "Keygrip = ") {
			continue
		}
		line = strings.TrimSpace(line)
		keyGrip := strings.Split(line, " = ")[1]
		fmt.Println(prefix, "|", "adding preset for", keyGrip)
		goexec.MustRunPrefixed(prefix, presetBin,
			"--passphrase", keyPassphrase,
			"--preset", keyGrip, // keygrip
		)
	}
}
