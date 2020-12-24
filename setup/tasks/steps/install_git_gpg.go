package steps

import (
	"runtime"
	"strings"

	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

// InstallGitGpg will install git & gpg using the systems package manager
func InstallGitGpg(password string) {
	prefix := colorchooser.Sprint("install-git-gpg")

	if runtime.GOOS == "windows" {
		InstallScoop()
		goexec.MustRunPrefixed(prefix, "powershell",
			"-Command", "scoop install git gpg",
		)
		goexec.MustRunPrefixed(prefix, "git", "config", "--global", "core.eol", "lf")
		goexec.MustRunPrefixed(prefix, "git", "config", "--global", "core.autocrlf", "false")
		goexec.MustRunPrefixed(prefix, "git", "config", "--global", "credential.helper", "manager")
		return
	}

	if runtime.GOOS == "linux" {
		if commandExists("dnf") {
			if isRoot() {
				goexec.MustRunPrefixed(prefix, "dnf", "install", "-y", "git", "gnupg")
			} else {
				goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd("sudo",
					goexec.SetIn(strings.NewReader(password)),
					goexec.Args("-S", "dnf", "install", "-y", "git", "gnupg"),
				))
			}
			return
		}
		if commandExists("apt") {
			script := "apt update && apt install -y git gnupg"
			if isRoot() {
				goexec.MustRunPrefixed(prefix, "bash", "-c", script)
			} else {
				goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd("sudo",
					goexec.SetIn(strings.NewReader(password)),
					goexec.Args("-S", "bash", "-c", script),
				))
			}
			return
		}
	}

	goerr.Check(UnSupportedOsError)
}
