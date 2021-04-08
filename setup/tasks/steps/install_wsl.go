// +build windows

package steps

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/brad-jones/dotfiles/setup/tasks/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/cavaliercoder/grab"
	"github.com/hpcloud/tail"
)

// https://docs.microsoft.com/en-us/windows/wsl/install-win10

func MustInstallWSL(replace bool) {
	prefix := colorchooser.Sprint("install-wsl")

	if utils.CommandExists("wsl") {
		if !replace {
			fmt.Println(prefix, "|", "wsl is already installed")
			return
		}

		goexec.MustRunPrefixed(prefix, "wsl", "--shutdown")
	}

	tempDir, err := ioutil.TempDir("", "bradsDotFiles")
	goerr.Check(err)
	defer os.RemoveAll(tempDir)
	defer fmt.Println(prefix, "|", "deleted", tempDir)
	fmt.Println(prefix, "|", "created", tempDir)

	installer := filepath.Join(tempDir, "wsl_update_x64.msi")
	downloader := task.New(func() {
		fmt.Println(prefix, "|", "downloading wsl_update_x64.msi")
		_, err = grab.Get(installer, "https://wslstorestorage.blob.core.windows.net/wslblob/wsl_update_x64.msi")
		goerr.Check(err)
	})

	fmt.Println(prefix, "|", "enable feature Microsoft-Windows-Subsystem-Linux")
	goexec.MustRunPrefixed(prefix, utils.SudoBin(), "dism", "/online", "/enable-feature", "/featurename:Microsoft-Windows-Subsystem-Linux", "/all", "/norestart")

	fmt.Println(prefix, "|", "enable feature VirtualMachinePlatform")
	goexec.MustRunPrefixed(prefix, utils.SudoBin(), "dism", "/online", "/enable-feature", "/featurename:VirtualMachinePlatform", "/all", "/norestart")

	// DOCS say to reboot here, not sure how to automate that...

	downloader.MustWait()

	// TODO: hash check

	fmt.Println(prefix, "|", "installing wsl_update_x64.msi")
	logfile := filepath.Join(tempDir, "logs.txt")
	goexec.MustRunPrefixed(prefix, utils.SudoBin(), installer, "/quiet", "/norestart", "/log", logfile)

	t, err := tail.TailFile(logfile, tail.Config{Follow: true, Poll: true})
	goerr.Check(err)
	go t.StopAtEOF()
	for line := range t.Lines {
		fmt.Println(prefix, "|", line.Text)
	}

	fmt.Println(prefix, "|", "set WSL2 as default engine")
	goexec.MustRunPrefixed(prefix, "wsl", "--set-default-version", "2")
}

func InstallWSLAsync(replace bool) *task.Task {
	return task.New(func() { MustInstallWSL(replace) })
}
