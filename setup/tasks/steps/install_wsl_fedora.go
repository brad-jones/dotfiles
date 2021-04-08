// +build windows

package steps

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/brad-jones/dotfiles/setup/tasks/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/cavaliercoder/grab"
	"github.com/gosimple/slug"
	"github.com/mholt/archiver"
)

// https://github.com/yosukes-dev/FedoraWSL

func MustInstallWSLFedora(version string, makeDefault, replace bool) string {
	wslVmName := fmt.Sprintf("Fedora%s", strings.Split(version, ".")[0])
	prefix := colorchooser.Sprint("install-wsl-" + slug.Make(wslVmName))

	extracted := filepath.Join(utils.HomeDir(), ".wsl", wslVmName)
	if utils.FolderExists(extracted) {
		fmt.Println(prefix, "|", "distro already installed")
		if !replace {
			return wslVmName
		}

		fmt.Println(prefix, "|", "replacing distro")
		goexec.RunPrefixed(prefix, "wsl", "--unregister", wslVmName)
		os.RemoveAll(extracted)
	}

	tempDir, err := ioutil.TempDir("", "bradsDotFiles")
	goerr.Check(err)
	defer os.RemoveAll(tempDir)
	defer fmt.Println(prefix, "|", "deleted", tempDir)
	fmt.Println(prefix, "|", "created", tempDir)

	fmt.Println(prefix, "|", "downloading Fedora WSL distro")
	zip := filepath.Join(tempDir, "fedora.zip")
	_, err = grab.Get(zip, fmt.Sprintf("https://github.com/yosukes-dev/FedoraWSL/releases/download/%s/%s.zip", version, wslVmName))
	goerr.Check(err)

	fmt.Println(prefix, "|", "extracting", zip)
	goerr.Check(archiver.Unarchive(zip, extracted), zip, extracted)

	goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd(
		filepath.Join(extracted, wslVmName+".exe"),
		goexec.SetIn(strings.NewReader("\n")),
	))

	if makeDefault {
		fmt.Println(prefix, "|", "making default")
		goexec.MustRunPrefixed(prefix, "wsl", "--set-default", wslVmName)
	}

	return wslVmName
}

func InstallWSLFedoraAsync(version string, makeDefault, replace bool) *task.Task {
	return task.New(func() { MustInstallWSLFedora(version, makeDefault, replace) })
}
