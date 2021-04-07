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

func MustInstallWavebox() {
	prefix := colorchooser.Sprint("install-wavebox")

	if runtime.GOOS == "windows" {
		if fileExists(`C:\Users\brad.jones\AppData\Local\WaveboxApp\Application\wavebox.exe`) {
			fmt.Println(prefix, "|", "wavebox is already installed")
			return
		}

		downloadDir, err := ioutil.TempDir("", "bradsDotFiles")
		goerr.Check(err)
		defer os.RemoveAll(downloadDir)
		defer fmt.Println(prefix, "|", "deleted", downloadDir)
		fmt.Println(prefix, "|", "created", downloadDir)

		fmt.Println(prefix, "|", "downloading wavebox_installer.exe")
		installer := filepath.Join(downloadDir, "wavebox_installer.exe")
		_, err = grab.Get(installer, "https://download.wavebox.app/latest/stable/win")
		goerr.Check(err)

		fmt.Println(prefix, "|", "running wavebox_installer.exe")
		goexec.MustRunPrefixed(prefix, installer)

		for {
			time.Sleep(time.Millisecond * 100)

			procs, err := ps.Processes()
			goerr.Check(err)

			exists := false
			for _, p := range procs {
				if p.Executable() == "wavebox_installer.exe" {
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

		fmt.Println(prefix, "|", "wavebox is installed")
	}
}

func InstallWaveboxAsync() *task.Task {
	return task.New(func() { MustInstallWavebox() })
}
