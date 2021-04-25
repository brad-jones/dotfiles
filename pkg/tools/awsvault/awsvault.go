package awsvault

import (
	"path/filepath"
	"runtime"

	"github.com/brad-jones/dotfiles/pkg/tools"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/dotfiles/pkg/utils/ghpkg"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
)

// https://github.com/99designs/aws-vault/releases/tag/v5.4.4

func Install() (err error) {
	defer goerr.Handle(func(e error) { err = e })
	v := tools.GetVersion("aws-vault")
	pkgPattern := "aws-vault-" + runtime.GOOS + "-"
	if runtime.GOOS == "linux" {
		pkgPattern = pkgPattern + runtime.GOARCH
	}
	if runtime.GOOS == "windows" {
		pkgPattern = pkgPattern + "386.exe"
	}
	ghpkg.MustInstallPkg("99designs", "aws-vault", v.No,
		ghpkg.Sha256Hash(v.Hash),
		ghpkg.PkgPattern(pkgPattern),
		ghpkg.ExeName("aws-vault"),
		ghpkg.Naked(true),
	)
	return
}

func MustInstall() {
	goerr.Check(Install())
}

func InstallAsync() *task.Task {
	return task.New(func() { MustInstall() })
}

func Path() string {
	p := filepath.Join(utils.HomeDir(), ".local", "bin", "aws-vault")
	if runtime.GOOS == "windows" {
		p = p + ".exe"
	}
	return p
}
