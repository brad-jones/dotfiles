package steps

import (
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

// MustInstallGitGpg will install git & gpg using the systems package manager
func MustInstallGitGpg(password string) {
	prefix := colorchooser.Sprint("install-git-gpg")

	MustInstallScoop()
	goexec.MustRunPrefixed(prefix, "powershell",
		"-Command", "scoop install git gpg",
	)
	goexec.MustRunPrefixed(prefix, "git", "config", "--global", "core.eol", "lf")
	goexec.MustRunPrefixed(prefix, "git", "config", "--global", "core.autocrlf", "false")
	goexec.MustRunPrefixed(prefix, "git", "config", "--global", "credential.helper", "manager")
}

func InstallGitGpgAsync(password string) *task.Task {
	return task.New(func() { MustInstallGitGpg(password) })
}
