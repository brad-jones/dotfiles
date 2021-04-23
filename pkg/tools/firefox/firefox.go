package firefox

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/brad-jones/dotfiles/pkg/tools/winsudo"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/cavaliercoder/grab"
	"github.com/shirou/gopsutil/v3/process"
)

func MustInstall() {
	prefix := colorchooser.Sprint("install-firefox")

	if runtime.GOOS == "windows" {
		if utils.FileExists(`C:\Program Files\Mozilla Firefox\firefox.exe`) {
			fmt.Println(prefix, "|", "skipping, already installed")
			return
		}

		downloadDir, err := ioutil.TempDir("", "bradsDotFiles")
		goerr.Check(err)
		defer os.RemoveAll(downloadDir)
		defer fmt.Println(prefix, "|", "deleted", downloadDir)
		fmt.Println(prefix, "|", "created", downloadDir)

		fmt.Println(prefix, "|", "downloading firefox_installer.exe")
		installer := filepath.Join(downloadDir, "firefox_installer.exe")
		_, err = grab.Get(installer, "https://download.mozilla.org/?product=firefox-latest&os=win64&lang=en-US")
		goerr.Check(err)

		fmt.Println(prefix, "|", "running firefox_installer.exe")
		goexec.MustRunPrefixed(prefix, winsudo.Path(), installer, "/S")

		// this is only needed with gsudo instead of my winsudo package
		for {
			time.Sleep(time.Millisecond * 100)

			procs, err := process.Processes()
			goerr.Check(err, "failed to list processes")

			exists := false
			for _, p := range procs {
				pName, err := p.Name()
				goerr.Check(err, "failed to get proc name")

				if pName == "firefox_installer.exe" {
					exists = true
					break
				}
			}

			if !exists {
				break
			}
		}

		// automatic installation of extensions & other settings will be done
		// via the browsers own "cloud sync" functionality.
		// We may be able to further improve this setup with something like:
		// https://github.com/go-vgo/robotgo

		fmt.Println(prefix, "|", "firefox is installed")
	}
}

func InstallAsync() *task.Task {
	return task.New(func() { MustInstall() })
}
