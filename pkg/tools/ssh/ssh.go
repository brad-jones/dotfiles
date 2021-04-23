package ssh

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/ActiveState/termtest/expect"
	"github.com/brad-jones/dotfiles/pkg/tools/scoop"
	"github.com/brad-jones/dotfiles/pkg/tools/winsudo"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/brad-jones/gopwsh"
	"github.com/gosimple/slug"
)

var BadPassPhrase = goerr.New("bad passphrase for ssh key")

var CouldNotKill = goerr.New("ssh-add not killed")

// UnSupportedOS is returned when we tried to install on an OS we don't recognize
var UnSupportedOS = goerr.New("we do not support the OS you are using")

// Installs the ssh tool
func Install(sudoPassword string) (err error) {
	defer goerr.Handle(func(e error) { err = e })
	prefix := colorchooser.Sprint("install-ssh")

	if runtime.GOOS == "windows" {
		goerr.Check(scoop.InstallOrUpdatePkgs(map[string]string{"win32-openssh": "*", "putty": "*"}), "failed to install ssh")
		return
	}

	if utils.CommandExists("dnf") {
		utils.RunElevatedNix(prefix, sudoPassword, "dnf", "install", "-y",
			"openssh",
			"openssh-clients",
		)

		// Solves: Bad owner or permissions on /etc/crypto-policies/back-ends/openssh.config
		//
		// I haven't been able to find any defintivie answers on this one,
		// again I think it's related to being a container / WSL.
		//
		// Closest thing I could find...
		// see: https://bugzilla.redhat.com/show_bug.cgi?id=1902646
		if utils.IsWSL() {
			fmt.Println(prefix, "|", "fix ssh permissions issue")
			utils.RunElevatedNix(prefix, sudoPassword,
				"chmod", "0644", "/etc/crypto-policies/back-ends/openssh.config",
			)
		}

		return
	}

	goerr.Check(UnSupportedOS)
	return
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
	prefix := colorchooser.Sprint("ssh-agent")
	fmt.Println(prefix, "|", "starting...")

	if runtime.GOOS == "windows" {
		ps := gopwsh.MustNew(gopwsh.Elevated(winsudo.Path()))
		defer ps.Exit()
		ps.MustExecute("Start-Service ssh-agent")
	}

	if utils.IsWSL() {
		utils.KillProcByName("ssh-agent")
		sock := filepath.Join(utils.HomeDir(), ".ssh/agent.sock")
		goexec.MustRunBuffered("ssh-agent", "-a", sock)
		goerr.Check(os.Setenv("SSH_AUTH_SOCK", sock), "failed to set SSH_AUTH_SOCK", sock)
	}

	fmt.Println(prefix, "|", "running")
}

func MustUnlockKey(keyFilePath, keyPassphrase string) {
	prefix := colorchooser.Sprint("unlock-ssh-key-" + slug.Make(filepath.Base(keyFilePath)))
	fmt.Println(prefix, "|", "adding", keyFilePath, "to ssh-agent")

	var b bytes.Buffer
	c, err := expect.NewConsole(expect.WithStdout(&b))
	goerr.Check(err, "expect.NewConsole failed")
	defer c.Close()
	go c.ExpectEOF()

	cmd := goexec.MustCmd("ssh-add",
		goexec.SetIn(c.Tty()),
		goexec.SetOut(c.Tty()),
		goexec.SetErr(c.Tty()),
		goexec.Args(keyFilePath),
	)

	defer func() {
		if cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
			if err := cmd.Process.Kill(); err != nil {
				if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
					goerr.Check(CouldNotKill)
				}
			}
		}
	}()

	goerr.Check(c.Pty.StartProcessInTerminal(cmd))

	sentPass := false

	for {
		time.Sleep(time.Microsecond)
		stdout := b.String()

		if !sentPass && strings.Contains(stdout, "Enter passphrase for") {
			_, err = c.SendLine(keyPassphrase)
			goerr.Check(err)
			sentPass = true
		}

		if sentPass && strings.Contains(stdout, "Identity added") {
			break
		}

		if sentPass && strings.Contains(stdout, "Bad passphrase") {
			goerr.Check(BadPassPhrase)
		}
	}

	goerr.Check(cmd.Wait(), "failed waiting for ssh-add")
	fmt.Println(prefix, "|", "added", keyFilePath, "to ssh-agent")
}

func MustUnlockPagentKey(keyFilePath string) {
	prefix := colorchooser.Sprint("unlock-ssh-key-" + slug.Make(filepath.Base(keyFilePath)))

	defer os.Remove(keyFilePath)
	defer fmt.Println(prefix, "|", "deleting", keyFilePath)

	fmt.Println(prefix, "|", "adding", keyFilePath, "to pagent")

	// hacky way of demonizing pagent
	goexec.RunBufferedAsync("pageant", keyFilePath)
	time.Sleep(time.Second * 1)

	fmt.Println(prefix, "|", "added", keyFilePath, "to pagent")
}
