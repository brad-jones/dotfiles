package steps

import (
	"fmt"
	"path/filepath"

	"github.com/brad-jones/dotfiles/pkg/survey"
	"github.com/brad-jones/dotfiles/pkg/tools/dart"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

func MustInstallDartScripts(answers *survey.Answers) {
	prefix := colorchooser.Sprint("dart-scripts")
	fmt.Println(prefix, "|", "restoring deps")
	scriptsDir := filepath.Join(utils.HomeDir(), ".local", "sbin")
	goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd(dart.PubPath(), goexec.Args("get"), goexec.Cwd(scriptsDir)))
}

func InstallDartScriptsAsync(answers *survey.Answers) *task.Task {
	return task.New(func() { MustInstallDartScripts(answers) })
}
