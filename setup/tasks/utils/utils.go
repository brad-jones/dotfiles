package utils

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

func GetComputerName() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("COMPUTERNAME")
	}

	if IsWSL() {
		return os.Getenv("NAME")
	}

	return os.Getenv("HOSTNAME")
}

func IsWSL() bool {
	return strings.Contains(goexec.MustRunBuffered("uname", "-a").StdOut, "microsoft-standard-WSL2")
}

func KillProcByName(name string) {
	if runtime.GOOS == "windows" {
		ps := gopwsh.MustNew(gopwsh.Elevated(SudoBin()))
		defer ps.Exit()
		ps.MustExecute(fmt.Sprintf("Stop-Process -Name %s -Force", gopwsh.QuoteArg(name)))
	}
}

func IsRoot() bool {
	return goexec.MustRunBuffered("id", "-u").StdOut == "0"
}

func HomeDir() string {
	home, err := os.UserHomeDir()
	goerr.Check(err, "failed to get the users home dir")
	return home
}

func ScoopDir() string {
	return os.Getenv("SCOOP")
}

func SudoBin() string {
	return filepath.Join(HomeDir(), ".local", "bin", "sudo.exe")
}

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FolderExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func SetWritable(filepath string) {
	goerr.Check(os.Chmod(filepath, 0222))
}

// this is super annoying, gpg writes to /dev/tty, instead of the usual
// /dev/stdout|stderr. In some cases "--batch" is meant to solve that but
// not for "--edit-key" it seems. So because I can't stream the output through
// my prefixer nor can I discard it, it corrupts the terminal output.
func TrustGpgKey(keyName string) {
	goexec.MustRunBufferedCmd(goexec.MustCmd("gpg",
		goexec.SetIn(strings.NewReader("5\ny\n")),
		goexec.Args(
			"--command-fd", "0",
			"--edit-key", keyName,
			"trust",
		),
	))
}

func ImportGpgKey(prefix, keyPath, keyName, keyPassphrase string) {
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
	TrustGpgKey("Brad Jones (vault) <brad@bjc.id.au>")

	fmt.Println(prefix, "|", "trusted", keyName)
}
