package steps

import (
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

// InstallScoopPkg will install the given scoop package
func InstallScoopPkg(pkgName, version string) {
	prefix := colorchooser.Sprint("install-scoop-" + pkgName)
	goexec.MustRunPrefixed(prefix, "powershell", "-Command", "scoop install "+pkgName+"@"+version)
}
