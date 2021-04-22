package assets

import (
	"runtime"

	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
)

// WriteDotfiles outputs all our home directory files that haven't been taken
// care of by other my specific tool installers. For example the gpg installer
// handles the writing of the gpg-agent.conf file.
//
// TODO: In actual fact I think I want to reduce this down to nothing.
// ie: We could end up with an aws-cli tool, a git tool, a vscode tool
// & a windows terminal tool.
func WriteDotfiles() (err error) {
	defer goerr.Handle(func(e error) { err = e })

	WriteFolderToHome(".aws")
	WriteFolderToHome("Projects")
	WriteFileToHome(".gitconfig")

	if runtime.GOOS == "windows" {
		WriteFolderToHome("AppData/Local/Microsoft/Windows Terminal")
		WriteFolderToHome("AppData/Roaming/Code")
		WriteFolderToHome("Documents")
	}

	return
}

func MustWriteDotfiles() {
	goerr.Check(WriteDotfiles())
}

func WriteDotfilesAsync() *task.Task {
	return task.New(func() { MustWriteDotfiles() })
}
