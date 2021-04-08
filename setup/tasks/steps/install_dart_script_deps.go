package steps

import (
	"path/filepath"
	"runtime"

	"github.com/brad-jones/dotfiles/setup/tasks/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

func MustInstallDartScriptDeps() {
	prefix := colorchooser.Sprint("install-dart-script-deps")
	scriptsDir := filepath.Join(utils.HomeDir(), ".local", "sbin")

	pub := "pub"
	if runtime.GOOS == "windows" {
		pub = filepath.Join(utils.ScoopDir(), "apps", "dart", "current", "bin", "pub.bat")
		utils.SetWritable(filepath.Join(scriptsDir, "pubspec.lock"))
	}

	goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd(pub, goexec.Args("get"), goexec.Cwd(scriptsDir)))
}

func InstallDartScriptDepsAsync() *task.Task {
	return task.New(func() { MustInstallDartScriptDeps() })
}
