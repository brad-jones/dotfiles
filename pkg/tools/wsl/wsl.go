package wsl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/brad-jones/dotfiles/pkg/tools"
	"github.com/brad-jones/dotfiles/pkg/tools/winsudo"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/dotfiles/pkg/utils/downloader"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/hpcloud/tail"
)

const wslMSIUpdate = "https://wslstorestorage.blob.core.windows.net/wslblob/wsl_update_x64.msi"

var MsiHashInvalid = goerr.New("the downloaded MSI did not match it's expected hash")

// InstallWSL automates https://docs.microsoft.com/en-us/windows/wsl/install-win10
// At least until the new `wsl --install` functionality lands in a GA release of Windows.
func Install(reset bool) (err error) {
	defer goerr.Handle(func(e error) { err = e })
	prefix := colorchooser.Sprint("install-wsl")

	// Bail out early if the WSL CLI already exists.
	// Unless of course we have told to do a full reset.
	if utils.CommandExists("wsl") {
		if !reset {
			fmt.Println(prefix, "|", "wsl is already installed")
			return
		}

		// Make sure the WSL VM is shutdown before we re-install it
		goexec.MustRunPrefixed(prefix, "wsl", "--shutdown")
	}

	// Create a temp dir
	tempDir, err := ioutil.TempDir("", "bradsDotFiles")
	goerr.Check(err, "failed to create tmp dir")
	defer os.RemoveAll(tempDir)
	defer fmt.Println(prefix, "|", "deleted", tempDir)
	fmt.Println(prefix, "|", "created", tempDir)

	// Start downloading the wsl update MSI installer
	fmt.Println(prefix, "|", "downloading wsl_update_x64.msi")
	dl := downloader.DownloadWithProgressAsync(prefix, wslMSIUpdate, filepath.Join(tempDir, "."))

	// Enable the required Windows Features
	fmt.Println(prefix, "|", "enable feature Microsoft-Windows-Subsystem-Linux")
	goexec.MustRunPrefixed(prefix, winsudo.Path(),
		"dism", "/online", "/enable-feature",
		"/featurename:Microsoft-Windows-Subsystem-Linux",
		"/all", "/norestart",
	)

	fmt.Println(prefix, "|", "enable feature VirtualMachinePlatform")
	goexec.MustRunPrefixed(prefix, winsudo.Path(),
		"dism", "/online", "/enable-feature",
		"/featurename:VirtualMachinePlatform",
		"/all", "/norestart",
	)

	// TODO: DOCS say to reboot here, not sure how to automate that...

	// Wait for the msi download
	r, err := dl.Result()
	goerr.Check(err, "failed to download", wslMSIUpdate)
	installer := r.(string)

	// Check the hash
	if utils.Sha256HashFile(installer) != tools.GetVersion(wslMSIUpdate).Hash {
		goerr.Check(MsiHashInvalid)
	}

	// Install the MSI, this will run detached
	fmt.Println(prefix, "|", "installing wsl_update_x64.msi")
	logfile := filepath.Join(tempDir, "logs.txt")
	goexec.MustRunPrefixed(prefix, winsudo.Path(),
		installer, "/quiet", "/norestart", "/log", logfile,
	)

	// Tail the log file so we know when the installer finished
	t, err := tail.TailFile(logfile, tail.Config{Follow: true, Poll: true})
	goerr.Check(err, "failed to tail log file", logfile)
	go t.StopAtEOF()
	for line := range t.Lines {
		fmt.Println(prefix, "|", line.Text)
	}

	// Make sure we are using WSL2 and not WSL1
	fmt.Println(prefix, "|", "set WSL2 as default engine")
	goexec.MustRunPrefixed(prefix, "wsl", "--set-default-version", "2")
	return
}

func MustInstall(reset bool) {
	goerr.Check(Install(reset))
}

func InstallAsync(reset bool) *task.Task {
	return task.New(func() { MustInstall(reset) })
}
