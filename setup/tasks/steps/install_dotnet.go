package steps

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/brad-jones/dotfiles/setup/tasks/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/cavaliercoder/grab"
	"github.com/gosimple/slug"
)

func MustInstallDotnet(versions ...string) {
	prefix := colorchooser.Sprint("install-dotnet")

	downloadDir, err := ioutil.TempDir("", "bradsDotFiles")
	goerr.Check(err)
	defer os.RemoveAll(downloadDir)
	defer fmt.Println(prefix, "|", "deleted", downloadDir)
	fmt.Println(prefix, "|", "created", downloadDir)

	scriptDst := filepath.Join(downloadDir, "dotnet-install")
	scriptLink := "https://dot.net/v1/dotnet-install"
	if runtime.GOOS == "windows" {
		scriptLink = scriptLink + ".ps1"
		scriptDst = scriptDst + ".ps1"
	} else {
		scriptLink = scriptLink + ".sh"
		scriptDst = scriptDst + ".sh"
	}

	fmt.Println(prefix, "|", "downloading", scriptLink)
	_, err = grab.Get(scriptDst, scriptLink)
	goerr.Check(err)

	// TODO: hash checking

	installDir := filepath.Join(utils.HomeDir(), ".dotnet")

	for _, version := range versions {
		prefix = colorchooser.Sprint("install-dotnet-" + slug.Make(version))

		if runtime.GOOS == "windows" {
			goexec.MustRunPrefixed(prefix, "powershell", scriptDst,
				"-InstallDir", installDir,
				"-Channel", "Current",
				"-Version", version,
				"-NoPath",
			)
			continue
		}

		goexec.MustRunPrefixed(prefix, "bash", scriptDst,
			"--install-dir", installDir,
			"--channel", "Current",
			"--version", version,
			"--no-path",
		)
	}
}

func InstallDotnetAsync(versions ...string) *task.Task {
	return task.New(func() { MustInstallDotnet(versions...) })
}
