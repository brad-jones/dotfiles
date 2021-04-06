package steps

import (
	"path/filepath"

	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

func MustInstallDartScriptDeps() {
	prefix := colorchooser.Sprint("install-dart-script-deps")
	pub := filepath.Join(scoopDir(), "apps", "dart", "current", "bin", "pub.bat")
	scriptsDir := filepath.Join(homeDir(), ".local", "sbin")
	goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd(pub, goexec.Args("get"), goexec.Cwd(scriptsDir)))
}

func InstallDartScriptDepsAsync() *task.Task {
	return task.New(func() { MustInstallDartScriptDeps() })
}
