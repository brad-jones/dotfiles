package steps

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/brad-jones/goasync/v2/await"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

// InstallSSHGpgKeys will get my identity keys (as opposed to the gopass vault
// encryption key) out of the gopass vault and install them into SSH & GPG agents.
//
// TODO: I want to replace with some sort of password manager that
// can also act as an SSH & GPG agent, that way the keys always remain
// inside the vault.
func InstallSSHGpgKeys() {
	prefix := colorchooser.Sprint("install-ssh-gpg-keys")

	homeDir, err := os.UserHomeDir()
	goerr.Check(err)
	sshDir := filepath.Join(homeDir, ".ssh")
	fmt.Println(prefix, "creating", sshDir, "(if not found)")
	goerr.Check(os.MkdirAll(filepath.Join(homeDir, ".ssh"), 0644))

	tmpDir, err := ioutil.TempDir("", "bradsDotFiles")
	goerr.Check(err)
	defer func() {
		fmt.Println(prefix, "deleting", tmpDir)
		goerr.Check(os.RemoveAll(tmpDir))
	}()

	tasks := []*task.Task{
		task.New(func() {
			goexec.MustRunPrefixed(prefix, "gopass", "bin", "cp",
				"keys/ssh/brad@bjc.id.au",
				filepath.Join(sshDir, "brad@bjc.id.au"),
			)
		}),
		task.New(func() {
			dst := filepath.Join(tmpDir, "brad@bjc.id.au")
			goexec.MustRunPrefixed(prefix, "gopass", "bin", "cp",
				"keys/gpg/brad@bjc.id.au", dst,
			)
			goexec.MustRunPrefixed(prefix, "gpg", "--import", dst)
			trustGpgKey(prefix, "Brad Jones <brad@bjc.id.au>")
		}),
	}

	if os.Getenv("COMPUTERNAME") == "XLW-5CD936CWNQ" {
		tasks = append(tasks,
			task.New(func() {
				goexec.MustRunPrefixed(prefix, "gopass", "bin", "cp",
					"keys/ssh/brad.jones@xero.com",
					filepath.Join(sshDir, "brad.jones@xero.com"),
				)
			}),
			task.New(func() {
				dst := filepath.Join(tmpDir, "brad.jones@xero.com")
				goexec.MustRunPrefixed(prefix, "gopass", "bin", "cp",
					"keys/gpg/brad.jones@xero.com", dst,
				)
				goexec.MustRunPrefixed(prefix, "gpg", "--import", dst)
				trustGpgKey(prefix, "Brad Jones <brad.jones@xero.com>")
			}),
		)
	}

	await.MustFastAllOrError(tasks...)
}
