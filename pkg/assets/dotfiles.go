package assets

import (
	"path/filepath"
	"runtime"

	"github.com/brad-jones/dotfiles/pkg/tools/scoop"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
)

// WriteDotfiles outputs all our home directory files.
func WriteDotfiles(withSecrets bool) (err error) {
	defer goerr.Handle(func(e error) { err = e })

	if !withSecrets {
		WriteFolderToHome(".ssh")
		WriteFolderToHome(".config/gopass")

		if runtime.GOOS != "windows" {
			WriteFolderToHome(".gnupg")
		} else {
			WriteFolder(".gnupg",
				filepath.Join(scoop.Path(), "apps/gpg/current/home"),
			)
		}
		return
	}

	WriteFolderToHome(".aws")
	WriteFolderToHome(".local")
	WriteFolderToHome("Projects")
	WriteFileToHome(".gitconfig.tmpl")

	if runtime.GOOS == "windows" {
		WriteFolderToHome("AppData")
		WriteFolderToHome("Documents")
	} else {
		WriteFileToHome(".bashrc.tmpl")
		WriteFileToHome(".bash_profile")
		WriteFileToHome(".bash_logout")
	}

	return
}

func MustWriteDotfiles(withSecrets bool) {
	goerr.Check(WriteDotfiles(withSecrets))
}

func WriteDotfilesAsync(withSecrets bool) *task.Task {
	return task.New(func() { MustWriteDotfiles(withSecrets) })
}
