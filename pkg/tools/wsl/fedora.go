package wsl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/dotfiles/pkg/utils/downloader"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/gosimple/slug"
	"github.com/mholt/archiver"
)

var FedoraWSLHashInvalid = goerr.New("the downloaded WSL setup exe did not match it's expected hash")

// InstallFedora creates an instance of https://github.com/yosukes-dev/FedoraWSL
func InstallFedora(version, hash string, makeDefault, reset bool) (name string, err error) {
	defer goerr.Handle(func(e error) { err = e })

	name = fmt.Sprintf("Fedora%s", strings.Split(version, ".")[0])
	prefix := colorchooser.Sprint("install-wsl-" + slug.Make(name))
	extracted := filepath.Join(utils.HomeDir(), ".wsl", name)

	// If the ~/.wsl/name folder already exists lets just assume it has
	// already been installed and bail out. Unless we have been told to do
	// a reset in which case we will deregister the instance.
	if utils.FolderExists(extracted) {
		fmt.Println(prefix, "|", "distro already installed")
		if !reset {
			return
		}
		fmt.Println(prefix, "|", "replacing distro")
		goexec.RunPrefixed(prefix, "wsl", "--unregister", name)
		goerr.Check(os.RemoveAll(extracted), "failed to delete", extracted)
	}

	// Create a tmp dir
	tempDir, err := ioutil.TempDir("", "bradsDotFiles")
	goerr.Check(err, "failed to create tmp dir")
	defer os.RemoveAll(tempDir)
	defer fmt.Println(prefix, "|", "deleted", tempDir)
	fmt.Println(prefix, "|", "created", tempDir)

	// Download the Fedora WSL release
	fmt.Println(prefix, "|", "downloading Fedora WSL distro")
	url := fmt.Sprintf("https://github.com/yosukes-dev/FedoraWSL/releases/download/%s/%s.zip", version, name)
	zip := downloader.MustDownloadWithProgress(prefix, url, filepath.Join(tempDir, "."))

	// Check the downloads hash
	if utils.Sha256HashFile(zip) != hash {
		goerr.Check(FedoraWSLHashInvalid)
	}

	// Extract the WSL release archive
	fmt.Println(prefix, "|", "extracting", zip)
	goerr.Check(archiver.Unarchive(zip, extracted), zip, extracted)

	// Run the WSL setup
	goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd(
		filepath.Join(extracted, name+".exe"),
		goexec.SetIn(strings.NewReader("\n")),
	))

	// Set it as the default WSL instance
	if makeDefault {
		fmt.Println(prefix, "|", "making default")
		goexec.MustRunPrefixed(prefix, "wsl", "--set-default", name)
	}

	return
}

func MustInstallFedora(version, hash string, makeDefault, reset bool) string {
	name, err := InstallFedora(version, hash, makeDefault, reset)
	goerr.Check(err)
	return name
}

func InstallFedoraAsync(version, hash string, makeDefault, reset bool) *task.Task {
	return task.New(func() { MustInstallFedora(version, hash, makeDefault, reset) })
}
