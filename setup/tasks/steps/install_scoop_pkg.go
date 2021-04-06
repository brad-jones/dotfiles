package steps

import (
	"fmt"
	"path/filepath"

	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/brad-jones/gopwsh"
)

// MustInstallScoopPkg will install the given scoop package
func MustInstallScoopPkg(pkgName, version string) {
	prefix := colorchooser.Sprint("install-scoop-" + pkgName)

	if len(version) > 0 {
		goexec.MustRunPrefixed(prefix, "powershell", "-Command", "scoop install "+pkgName+"@"+version)
		return
	}

	if fileExists(filepath.Join(homeDir(), "scoop", "apps", pkgName, "current", "manifest.json")) {
		goexec.MustRunPrefixed(prefix, "powershell", "-Command", "scoop update "+pkgName)
		return
	}

	goexec.MustRunPrefixed(prefix, "powershell", "-Command", "scoop install "+pkgName)
}

func InstallScoopPkgAsync(pkgName, version string) *task.Task {
	return task.New(func() { MustInstallScoopPkg(pkgName, version) })
}

// MustInstallScoopPkgs will install the given scoop package
func MustInstallScoopPkgs(packages map[string]string) {
	prefix := colorchooser.Sprint("install-scoop-packages")

	ps := gopwsh.MustNew()
	defer ps.Exit()

	for pkg, ver := range packages {
		var stdout string
		if ver == "" && fileExists(filepath.Join(scoopDir(), "apps", pkg, "current", "manifest.json")) {
			stdout, _ = ps.MustExecute(fmt.Sprintf("scoop update %s", pkg))
		} else {
			stdout, _ = ps.MustExecute(fmt.Sprintf("scoop install %s@%s", pkg, ver))
		}
		fmt.Println(prefix, "|", stdout)
	}
}

func InstallScoopPkgsAsync(packages map[string]string) *task.Task {
	return task.New(func() { MustInstallScoopPkgs(packages) })
}
