package steps

import (
	"strings"

	"github.com/brad-jones/dotfiles/setup/tasks/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

// MustInstallGitGpg will install git & gpg using the systems package manager
func MustInstallGitGpg(password string) {
	prefix := colorchooser.Sprint("install-git-gpg")

	if utils.CommandExists("dnf") {
		if utils.IsRoot() {
			goexec.MustRunPrefixed(prefix, "dnf", "install", "-y", "git", "gnupg", "pinentry")
		} else {
			if len(password) > 0 {
				goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd("sudo",
					goexec.SetIn(strings.NewReader(password)),
					goexec.Args("-S", "dnf", "install", "-y", "git", "gnupg", "pinentry"),
				))
			} else {
				goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd("sudo",
					goexec.Args("dnf", "install", "-y", "git", "gnupg", "pinentry"),
				))
			}
		}
		return
	}

	if utils.CommandExists("apt") {
		script := "apt update && apt install -y git gnupg"
		if utils.IsRoot() {
			goexec.MustRunPrefixed(prefix, "bash", "-c", script)
		} else {
			if len(password) > 0 {
				goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd("sudo",
					goexec.SetIn(strings.NewReader(password)),
					goexec.Args("-S", "bash", "-c", script),
				))
			} else {
				goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd("sudo",
					goexec.Args("bash", "-c", script),
				))
			}
		}
		return
	}
}

func InstallGitGpgAsync(password string) *task.Task {
	return task.New(func() { MustInstallGitGpg(password) })
}
