package wsl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/brad-jones/dotfiles/pkg/assets"
	"github.com/brad-jones/dotfiles/pkg/survey"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/gosimple/slug"
)

// InstallDotfiles is a little recursive, after doing some initial prep work
// in the WSL instance it runs this same tool but inside the Linux WSL instance.
func InstallDotfiles(name string, answers *survey.Answers) (err error) {
	defer goerr.Handle(func(e error) { err = e })
	prefix := colorchooser.Sprint("wsl-" + slug.Make(name))

	// Extract the linux version of this tool into a temp location
	tempDir, err := ioutil.TempDir("", "bradsDotFiles")
	goerr.Check(err, "failed to create temp dir")
	defer os.RemoveAll(tempDir)
	defer fmt.Println(prefix, "|", "deleted", tempDir)
	fmt.Println(prefix, "|", "created", tempDir)
	linBin := filepath.Join(tempDir, "dotfiles")
	assets.WriteFile("dotfiles_linux_amd64", linBin)
	fmt.Println(prefix, "|", "extracted", linBin)

	// Run ourselves inside the wsl instance
	linBin = strings.Replace(linBin, "C:\\", "/mnt/c/", 1)
	linBin = strings.ReplaceAll(linBin, "\\", "/")
	fmt.Println(prefix, "|", "running", linBin)
	j, err := json.Marshal(answers)
	goerr.Check(err, "failed to stringify answers")
	goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd("wsl",
		goexec.SetIn(bytes.NewReader(j)),
		goexec.Args("-d", name, "-e", linBin),
	))

	return
}

func MustInstallDotfiles(name string, answers *survey.Answers) {
	goerr.Check(InstallDotfiles(name, answers))
}

func InstallDotfilesAsync(name string, answers *survey.Answers) *task.Task {
	return task.New(func() { MustInstallDotfiles(name, answers) })
}
