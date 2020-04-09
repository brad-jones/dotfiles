package steps

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/brad-jones/dotfiles/pkg/assets"
	"github.com/brad-jones/dotfiles/pkg/tools/scoop"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

func MustInstallDartScripts() {
	prefix := colorchooser.Sprint("dart-scripts")

	if !utils.CommandExists("dart") {
		fmt.Println(prefix, "|", "installing dart")
		if runtime.GOOS == "windows" {
			scoop.MustInstallOrUpdatePkgs(map[string]string{"dart": "*"})
		}
		fmt.Println(prefix, "|", "dart installed")
	} else {
		fmt.Println(prefix, "|", "dart already installed")
	}

	fmt.Println(prefix, "|", "copying to filesystem")
	assets.WriteFolderToHome(".local/sbin")

	fmt.Println(prefix, "|", "restoring deps")
	scriptsDir := filepath.Join(utils.HomeDir(), ".local", "sbin")

	pub := "pub"
	if runtime.GOOS == "windows" {
		pub = filepath.Join(scoop.Path(), "apps", "dart", "current", "bin", "pub.bat")
	}

	goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd(pub, goexec.Args("get"), goexec.Cwd(scriptsDir)))
}

func InstallDartScriptsAsync() *task.Task {
	return task.New(func() { MustInstallDartScripts() })
}
