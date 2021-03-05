package steps

import (
	"os"
	"os/exec"
	"strings"

	"github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
)

func isRoot() bool {
	return goexec.MustRunBuffered("id", "-u").StdOut == "0"
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

func trustGpgKey(prefix, keyName string) {
	goexec.MustRunPrefixedCmd(prefix,
		goexec.MustCmd("gpg",
			goexec.SetIn(strings.NewReader("5\r\ny\r\n")),
			goexec.Args(
				"--command-fd", "0",
				"--edit-key", keyName,
				"trust",
			),
		),
	)
}

func isElevated() bool {
	ps, err := powershell.New(&backend.Local{})
	goerr.Check(err)
	defer ps.Exit()

	_, _, err = ps.Execute("$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent());")
	goerr.Check(err)
	stdout, _, err := ps.Execute("$currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator);")
	goerr.Check(err)

	return strings.TrimSpace(stdout) == "True"
}
