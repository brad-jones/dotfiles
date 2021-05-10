package steps

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/brad-jones/dotfiles/pkg/survey"
	"github.com/brad-jones/dotfiles/pkg/tools/gopass"
	"github.com/brad-jones/dotfiles/pkg/tools/gpg"
	"github.com/brad-jones/dotfiles/pkg/tools/ssh"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goasync/v2/await"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/gosimple/slug"
)

func UnlockKeys(answers *survey.Answers) (err error) {
	defer goerr.Handle(func(e error) { err = e })
	prefix := colorchooser.Sprint("unlock-keys")

	if len(answers.SudoPassword) == 0 {
		answers.SudoPassword = gopass.GetSecret(fmt.Sprintf("sudo/%s", slug.Make(utils.GetComputerName())))
	}
	ssh.MustInstall(answers.SudoPassword)
	ssh.MustStartAgent()

	sshDir := filepath.Join(utils.HomeDir(), ".ssh")
	tmpDir, err := ioutil.TempDir("", "bradsDotFiles")
	goerr.Check(err, "failed to create temp dir")
	defer os.RemoveAll(tmpDir)
	defer fmt.Println(prefix, "|", "deleting", tmpDir)

	tasks := []*task.Task{
		task.New(func() {
			dst := filepath.Join(sshDir, "brad@bjc.id.au")
			gopass.WriteBinarySecret("keys/ssh/brad@bjc.id.au", dst, 0600)
			ssh.MustUnlockKey(dst, gopass.GetSecret("keys/ssh/brad@bjc.id.au.pass"))
		}),
		task.New(func() {
			dst := filepath.Join(tmpDir, "brad@bjc.id.au")
			gopass.WriteBinarySecret("keys/gpg/brad@bjc.id.au", dst, 0600)
			gpg.MustImportKey(dst, "Brad Jones <brad@bjc.id.au>", gopass.GetSecret("keys/gpg/brad@bjc.id.au.pass"), true)
		}),
	}

	if runtime.GOOS == "windows" {
		tasks = append(tasks, task.New(func() {
			dst := filepath.Join(sshDir, "brad@bjc.id.au.ppk")
			gopass.WriteBinarySecret("keys/ssh/brad@bjc.id.au.ppk", dst, 0600)
			ssh.MustUnlockPagentKey(dst)
		}))
	}

	if utils.GetComputerName() == "XLW-5CD936CWNQ" {
		tasks = append(tasks,
			task.New(func() {
				dst := filepath.Join(sshDir, "brad.jones@xero.com")
				gopass.WriteBinarySecret("keys/ssh/brad.jones@xero.com", dst, 0600)
				ssh.MustUnlockKey(dst, gopass.GetSecret("keys/ssh/brad.jones@xero.com.pass"))
			}),
			task.New(func() {
				dst := filepath.Join(tmpDir, "brad.jones@xero.com")
				gopass.WriteBinarySecret("keys/gpg/brad.jones@xero.com", dst, 0600)
				gpg.MustImportKey(dst, "Brad Jones <brad.jones@xero.com>", gopass.GetSecret("keys/gpg/brad.jones@xero.com.pass"), true)
			}),
		)

		if runtime.GOOS == "windows" {
			tasks = append(tasks, task.New(func() {
				dst := filepath.Join(sshDir, "brad.jones@xero.com.ppk")
				gopass.WriteBinarySecret("keys/ssh/brad.jones@xero.com.ppk", dst, 0600)
				ssh.MustUnlockPagentKey(dst)
			}))
		}
	}

	await.MustFastAllOrError(tasks...)

	// Now that all our SSH keys are added to the agent we
	// can then pull any changes to the gopass vault.
	goexec.MustRunPrefixed(prefix, gopass.Path(), "sync")

	return
}

func MustUnlockKeys(answers *survey.Answers) {
	goerr.Check(UnlockKeys(answers))
}

func UnlockKeysAsync(answers *survey.Answers) *task.Task {
	return task.New(func() { MustUnlockKeys(answers) })
}
