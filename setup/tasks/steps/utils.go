package steps

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/gopwsh"
)

func isRoot() bool {
	return goexec.MustRunBuffered("id", "-u").StdOut == "0"
}

func homeDir() string {
	home, err := os.UserHomeDir()
	goerr.Check(err, "failed to get the users home dir")
	return home
}

func sudoBin() string {
	return filepath.Join(homeDir(), ".local", "bin", "sudo.exe")
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// this is super annoying, gpg writes to /dev/tty, instead of the usual
// /dev/stdout|stderr. In some cases "--batch" is meant to solve that but
// not for "--edit-key" it seems. So because I can't stream the output through
// my prefixer nor can I discard it, it corrupts the terminal output.
func trustGpgKey(keyName string) {
	// On Windows I can start a hidden console. ie: ShellExecute.
	// This is using my winsudo package, so it's running elevated when it
	// doesn't have to be. But if I call ShellExecute directly I then need a
	// way to catch errors, etc... so this will do for now.
	if runtime.GOOS == "windows" {
		goexec.MustRunBufferedCmd(goexec.MustCmd(sudoBin(),
			goexec.Args(
				"powershell",
				"-Command",
				"echo '5\r\ny\r\n' | gpg --command-fd 0 --edit-key "+gopwsh.QuoteArg(keyName)+" trust",
			),
		))
		return
	}

	// Else where lets just do this for now.
	goexec.MustRunBufferedCmd(goexec.MustCmd("gpg",
		goexec.SetIn(strings.NewReader("5\ny\n")),
		goexec.Args(
			"--command-fd", "0",
			"--edit-key", keyName,
			"trust",
		),
	))
}

func importGpgKey(prefix, keyPath, keyName string) {
	fmt.Println(prefix, "importing", keyPath)
	goexec.MustRunBuffered("gpg", "--import", keyPath)

	fmt.Println(prefix, "imported", keyPath, "trusting", keyName)
	trustGpgKey("Brad Jones (vault) <brad@bjc.id.au>")

	fmt.Println(prefix, "trusted", keyName)
}
