package scoop

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/brad-jones/dotfiles/pkg/tools/winsudo"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/brad-jones/gopwsh"
)

// NotWindows is returned whenever one of these tasks are run on a non-windows OS
var NotWindows = goerr.New("scoop is a windows only technology")

// Installs the https://scoop.sh/ package manager if not already installed otherwise updates scoop
func Install(reset, update bool) (err error) {
	defer goerr.Handle(func(e error) { err = e })
	if runtime.GOOS != "windows" {
		goerr.Check(NotWindows)
	}

	prefix := colorchooser.Sprint("install-scoop")

	if !reset && utils.CommandExists("scoop") {
		if update {
			goexec.MustRunPrefixed(
				colorchooser.Sprint("update-scoop"),
				"powershell", "-Command", "scoop update",
			)
		} else {
			fmt.Println(prefix, "|", "skipping, already installed")
		}
		return
	}

	fmt.Println(prefix, "|", "setting execution policy to RemoteSigned")
	psElevated := gopwsh.MustNew(gopwsh.Elevated(winsudo.Path()))
	defer psElevated.Exit()
	psElevated.MustExecute("Set-ExecutionPolicy RemoteSigned -scope CurrentUser")

	fmt.Println(prefix, "|", "installing into", Path())
	goexec.MustRunPrefixedCmd(prefix,
		goexec.MustCmd("powershell",
			goexec.EnvCombined(map[string]string{"SCOOP": Path()}),
			goexec.Args("-Command", "iwr -useb get.scoop.sh | iex"),
		),
	)

	/*
		Scoop requires git for some functions like adding new buckets.

		It also makes use of aria2 to speed up downloads, so may as well have
		it installed before installing the rest of the tools.

		see: https://github.com/lukesampson/scoop#multi-connection-downloads-with-aria2
	*/
	InstallOrUpdatePkgs(map[string]string{"git": "*", "aria2": "*"})
	return
}

// MustInstall does the same thing as Install but panics instead of returning an error
func MustInstall(reset, update bool) {
	goerr.Check(Install(reset, update))
}

// InstallAsync does the same thing as Install but asynchronously.
func InstallAsync(reset, update bool) *task.Task {
	return task.New(func() { MustInstall(reset, update) })
}

// Adds a new bucket to scoop
func AddBucket(bucketName, repo string) (err error) {
	defer goerr.Handle(func(e error) { err = e })
	if runtime.GOOS != "windows" {
		goerr.Check(NotWindows)
	}

	prefix := colorchooser.Sprint("add-scoop-bucket-" + bucketName)

	if BucketExists(bucketName) {
		fmt.Println(prefix, "|", "bucket already exists")
		return
	}

	if len(repo) > 0 {
		goexec.MustRunPrefixed(prefix, "powershell", "-Command",
			fmt.Sprintf("scoop bucket add %s %s",
				gopwsh.QuoteArg(bucketName),
				gopwsh.QuoteArg(repo),
			),
		)
		return
	}

	goexec.MustRunPrefixed(prefix, "powershell", "-Command",
		fmt.Sprintf("scoop bucket add %s",
			gopwsh.QuoteArg(bucketName),
		),
	)
	return
}

// MustAddBucket does the same thing as AddBucket but panics instead of returning an error
func MustAddBucket(bucketName, repo string) {
	goerr.Check(AddBucket(bucketName, repo))
}

// AddBucketAsync does the same thing as AddBucket but asynchronously.
func AddBucketAsync(bucketName, repo string) *task.Task {
	return task.New(func() { MustAddBucket(bucketName, repo) })
}

// Installs or updates a collection of packages.
// The map key is the package name & the value is a version number.
// Wildcard is accepted in which case it is simply omitted from scoop commands.
func InstallOrUpdatePkgs(packages map[string]string) (err error) {
	defer goerr.Handle(func(e error) { err = e })
	if runtime.GOOS != "windows" {
		goerr.Check(NotWindows)
	}

	prefix := colorchooser.Sprint("install-update-scoop-packages")
	ps := gopwsh.MustNew()
	defer ps.Exit()

	for pkg, ver := range packages {
		var stdout string

		if PkgExists(pkg, ver) {
			if ver != "*" {
				fmt.Println(prefix, "|", "skipping", pkg, ver, "already installed")
				continue
			}
			stdout, _ = ps.MustExecute(fmt.Sprintf("scoop update %s", pkg))
		} else {
			if ver == "*" {
				stdout, _ = ps.MustExecute(fmt.Sprintf("scoop install %s", pkg))
			} else {
				stdout, _ = ps.MustExecute(fmt.Sprintf("scoop install %s@%s", pkg, ver))
			}
		}

		for _, line := range strings.Split(stdout, "\n") {
			fmt.Println(prefix, "|", line)
		}
	}

	return
}

// MustInstallOrUpdatePkgs does the same thing as InstallOrUpdatePkgs but panics instead of returning an error
func MustInstallOrUpdatePkgs(packages map[string]string) {
	goerr.Check(InstallOrUpdatePkgs(packages))
}

// InstallOrUpdatePkgsAsync does the same thing as InstallOrUpdatePkgs but asynchronously.
func InstallOrUpdatePkgsAsync(packages map[string]string) *task.Task {
	return task.New(func() { MustInstallOrUpdatePkgs(packages) })
}

// PkgExists performs a quick sanity check against the filesystem,
// much faster than shelling out to PowerShell, to see if a package
// exists or not.
func PkgExists(name, version string) bool {
	if version == "*" {
		version = "current"
	}
	return utils.FileExists(
		filepath.Join(Path(), "apps", name, version, "manifest.json"),
	)
}

// BucketExists performs a quick sanity check against the filesystem,
// much faster than shelling out to PowerShell, to see if a bucket
// exists or not.
func BucketExists(name string) bool {
	return utils.FolderExists(
		filepath.Join(Path(), "buckets", name, ".git"),
	)
}

// Returns the path to the scoop directory.
// Honours the "SCOOP" environment variable if set,
// otherwise defaults "~/.scoop".
func Path() string {
	dir := os.Getenv("SCOOP")
	if dir == "" {
		dir = filepath.Join(utils.HomeDir(), ".scoop")
	}
	return dir
}
