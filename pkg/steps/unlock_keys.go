package steps

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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
)

func UnlockKeys(answers *survey.Answers) (err error) {
	defer goerr.Handle(func(e error) { err = e })
	prefix := colorchooser.Sprint("unlock-keys")

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
			goexec.MustRunPrefixed(prefix, "gopass", "bin", "cp", "keys/ssh/brad@bjc.id.au", dst)
			ssh.MustUnlockKey(dst, gopass.GetSecret("keys/ssh/brad@bjc.id.au.pass"))
		}),
		task.New(func() {
			dst := filepath.Join(sshDir, "brad@bjc.id.au.ppk")
			goexec.MustRunPrefixed(prefix, "gopass", "bin", "cp", "keys/ssh/brad@bjc.id.au.ppk", dst)
			ssh.MustUnlockPagentKey(dst)
		}),
		task.New(func() {
			dst := filepath.Join(tmpDir, "brad@bjc.id.au")
			goexec.MustRunPrefixed(prefix, "gopass", "bin", "cp", "keys/gpg/brad@bjc.id.au", dst)
			gpg.MustImportKey(dst, "Brad Jones <brad@bjc.id.au>", gopass.GetSecret("keys/gpg/brad@bjc.id.au.pass"), true)
		}),
	}

	if utils.GetComputerName() == "XLW-5CD936CWNQ" {
		tasks = append(tasks,
			task.New(func() {
				dst := filepath.Join(sshDir, "brad.jones@xero.com")
				goexec.MustRunPrefixed(prefix, "gopass", "bin", "cp", "keys/ssh/brad.jones@xero.com", dst)
				ssh.MustUnlockKey(dst, gopass.GetSecret("keys/ssh/brad.jones@xero.com.pass"))
			}),
			task.New(func() {
				dst := filepath.Join(sshDir, "brad.jones@xero.com.ppk")
				goexec.MustRunPrefixed(prefix, "gopass", "bin", "cp", "keys/ssh/brad.jones@xero.com.ppk", dst)
				ssh.MustUnlockPagentKey(dst)
			}),
			task.New(func() {
				dst := filepath.Join(tmpDir, "brad.jones@xero.com")
				goexec.MustRunPrefixed(prefix, "gopass", "bin", "cp", "keys/gpg/brad.jones@xero.com", dst)
				gpg.MustImportKey(dst, "Brad Jones <brad.jones@xero.com>", gopass.GetSecret("keys/gpg/brad.jones@xero.com.pass"), true)
			}),
		)
	}

	await.MustFastAllOrError(tasks...)

	// Now that all our SSH keys are added to the agent we
	// can then pull any changes to the gopass vault.
	goexec.MustRunPrefixed("gopass", "gopass", "sync")

	return
}

func MustUnlockKeys(answers *survey.Answers) {
	goerr.Check(UnlockKeys(answers))
}

func UnlockKeysAsync(answers *survey.Answers) *task.Task {
	return task.New(func() { MustUnlockKeys(answers) })
}

/*
	# Install additional SSH Keys
	# ------------------------------------------------------------------------------
	# TODO: Would be nice use gopass as an actual ssh-agent?
	if ($env:COMPUTERNAME -eq "XLW-5CD936CWNQ") {
		RmIfExists -Path $env:USERPROFILE/.ssh/keys;

		Exec -ScriptBlock { mkdir $env:USERPROFILE/.ssh/keys/xero-payroll-prod; }
		Exec -ScriptBlock { gopass bin cp keys/ssh/xero-payroll-prod/payroll-checkpoint.pem $env:USERPROFILE\.ssh\keys\xero-payroll-prod\payroll-checkpoint.pem; }
		Exec -ScriptBlock { gopass bin cp keys/ssh/xero-payroll-prod/payroll-dev-public.pem $env:USERPROFILE\.ssh\keys\xero-payroll-prod\payroll-dev-public.pem; }
		Exec -ScriptBlock { gopass bin cp keys/ssh/xero-payroll-prod/payroll-devops.pem $env:USERPROFILE\.ssh\keys\xero-payroll-prod\payroll-devops.pem; }

		Exec -ScriptBlock { mkdir $env:USERPROFILE/.ssh/keys/xero-payroll-test; }
		Exec -ScriptBlock { gopass bin cp keys/ssh/xero-payroll-test/payroll-checkpoint.pem $env:USERPROFILE\.ssh\keys\xero-payroll-test\payroll-checkpoint.pem; }
		Exec -ScriptBlock { gopass bin cp keys/ssh/xero-payroll-test/payroll-dev-public.pem $env:USERPROFILE\.ssh\keys\xero-payroll-test\payroll-dev-public.pem; }
		Exec -ScriptBlock { gopass bin cp keys/ssh/xero-payroll-test/payroll-devops.pem $env:USERPROFILE\.ssh\keys\xero-payroll-test\payroll-devops.pem; }

		Exec -ScriptBlock { mkdir $env:USERPROFILE/.ssh/keys/xero-payroll-uat; }
		Exec -ScriptBlock { gopass bin cp keys/ssh/xero-payroll-uat/payroll-checkpoint.pem $env:USERPROFILE\.ssh\keys\xero-payroll-uat\payroll-checkpoint.pem; }
		Exec -ScriptBlock { gopass bin cp keys/ssh/xero-payroll-uat/payroll-dev-public.pem $env:USERPROFILE\.ssh\keys\xero-payroll-uat\payroll-dev-public.pem; }
		Exec -ScriptBlock { gopass bin cp keys/ssh/xero-payroll-uat/payroll-devops.pem $env:USERPROFILE\.ssh\keys\xero-payroll-uat\payroll-devops.pem; }

		Exec -ScriptBlock { mkdir $env:USERPROFILE/.ssh/keys/xero-ps-paas-svc; }
		Exec -ScriptBlock { gopass bin cp keys/ssh/xero-ps-paas-svc/payroll-devops.pem $env:USERPROFILE\.ssh\keys\xero-ps-paas-svc\payroll-devops.pem; }
	}
*/
