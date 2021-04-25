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
	"github.com/shirou/gopsutil/v3/process"
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
	found := false
	if runtime.GOOS == "windows" && !strings.HasSuffix(name, ".exe") {
		name = name + ".exe"
	}
	fmt.Println(colorchooser.Sprint("kill-proc"), "|", "looking for", name)

	procs, err := process.Processes()
	goerr.Check(err, "failed to list processes")

	for _, p := range procs {
		pName, err := p.Name()
		if err != nil {
			if strings.Contains(err.Error(), "couldn't find pid") ||
				strings.Contains(err.Error(), "status: no such file or directory") {
				continue
			}
			goerr.Check(err, "failed to get proc name")
		}

		if pName == name {
			var killChildren func(p *process.Process)
			killChildren = func(p *process.Process) {
				pChildren, err := p.Children()
				if err != nil {
					if err.Error() != "process does not have children" {
						goerr.Check(err, "failed to get proc children", pName)
					}
					return
				}
				for _, v := range pChildren {
					killChildren(v)
					goerr.Check(v.Kill(), "failed to kill proc child", pName)
				}
			}
			killChildren(p)
			goerr.Check(p.Kill(), "failed to kill proc parent", pName)
			found = true
		}
	}

	if found {
		fmt.Println(colorchooser.Sprint("kill-proc"), "|", name, "has been slayed")
	} else {
		fmt.Println(colorchooser.Sprint("kill-proc"), "|", name, "not running")
	}
}

func RunElevatedNix(prefix, sudoPassword, cmd string, args ...string) {
	if runtime.GOOS == "windows" {
		goerr.Check(goerr.New("RunElevatedNix is not designed to work on Windows"))
	}

	if IsRoot() {
		goexec.MustRunPrefixed(prefix, cmd, args...)
		return
	}

	if len(sudoPassword) > 0 {
		goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd("sudo",
			goexec.SetIn(strings.NewReader(sudoPassword)),
			goexec.Args(append([]string{"-S", cmd}, args...)...),
		))
		return
	}

	goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd("sudo",
		goexec.Args(append([]string{cmd}, args...)...),
	))
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
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		goerr.Check(err, "failed to stat file", filename)
	}
	return !info.IsDir()
}

func FolderExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		goerr.Check(err, "failed to stat folder", filename)
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
