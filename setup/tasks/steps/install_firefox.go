package steps

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/cavaliercoder/grab"
	"github.com/mitchellh/go-ps"
)

func MustInstallFirefox() {
	prefix := colorchooser.Sprint("install-firefox")

	if runtime.GOOS == "windows" {
		if fileExists(`C:\Program Files\Mozilla Firefox\firefox.exe`) {
			fmt.Println(prefix, "|", "firefox is already installed")
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
		goexec.MustRunPrefixed(prefix, sudoBin(), installer, "/S")

		// this is only needed with gsudo instead of my winsudo package
		for {
			time.Sleep(time.Millisecond * 100)

			procs, err := ps.Processes()
			goerr.Check(err)

			exists := false
			for _, p := range procs {
				if p.Executable() == "firefox_installer.exe" {
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

func InstallFirefoxAsync() *task.Task {
	return task.New(func() { MustInstallFirefox() })
}
