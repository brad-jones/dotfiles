package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/brad-jones/gopwsh"
	"github.com/zalando/go-keyring"
)

func GetSecretFromKeychain(service, user string) string {
	secret, err := keyring.Get(service, user)
	goerr.Check(err, "failed to get secret from native OS keychain")
	return secret
}

func Sha256HashContent(dat []byte) string {
	h := sha256.New()
	h.Write(dat)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Sha256HashFile(filePath string) string {
	f, err := os.Open(filePath)
	goerr.Check(err, "failed to open", filePath)
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		goerr.Check(err, "failed to copy", filePath)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

// see: https://stackoverflow.com/questions/22744443
func IsATerminal() bool {
	stat, err := os.Stdin.Stat()
	goerr.Check(err, "IsATerminal failed to stat STDIN")
	return (stat.Mode() & os.ModeCharDevice) != 0
}

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
	if runtime.GOOS == "linux" {
		return strings.Contains(
			goexec.MustRunBuffered("uname", "-a").StdOut,
			"microsoft-standard-WSL2",
		)
	}
	return false
}

func KillProcByName(name string) {
	if runtime.GOOS == "windows" {
		ps := gopwsh.MustNew(gopwsh.Elevated())
		defer ps.Exit()
		ps.MustExecute(fmt.Sprintf("Stop-Process -Name %s -Force", gopwsh.QuoteArg(name)))
	}
	fmt.Println(colorchooser.Sprint("kill-proc"), "|", name)
}

func IsRoot() bool {
	return goexec.MustRunBuffered("id", "-u").StdOut == "0"
}

func HomeDir() string {
	home, err := os.UserHomeDir()
	goerr.Check(err, "failed to get the users home dir")
	return home
}

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err == nil {
		return true
	}

	if runtime.GOOS == "windows" {
		ps := gopwsh.MustNew()
		defer ps.Exit()
		if _, stderr := ps.MustExecute("Get-Command " + gopwsh.QuoteArg(cmd)); len(stderr) == 0 {
			return true
		}
	}

	return false
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

func FileContents(filename string) string {
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(out)
}
