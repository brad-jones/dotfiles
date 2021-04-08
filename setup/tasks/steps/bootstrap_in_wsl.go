// +build windows

package steps

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/brad-jones/dotfiles/setup/tasks/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/gosimple/slug"
)

//go:embed setup_linux_amd64
var linuxSetupBin []byte

func MustBoostrapInWSL(distro string) {
	prefix := colorchooser.Sprint("bootstrap-in-wsl-" + slug.Make(distro))
	distroCliTool := filepath.Join(utils.HomeDir(), ".wsl", distro, distro+".exe")

	currentUser := os.Getenv("USERNAME")
	currentWSLUser := strings.TrimSpace(goexec.MustRunBuffered("wsl", "-d", distro, "echo", "$USER").StdOut)
	if currentUser != currentWSLUser {
		fmt.Println(prefix, "|", "installing sudo")
		goexec.MustRunPrefixed(prefix, "wsl", "-d", distro, "dnf", "install", "-y", "sudo")

		fmt.Println(prefix, "|", "allow sudo without password")
		goexec.MustRunPrefixed(prefix, "wsl", "-d", distro, "sh", "-c", "echo '%wheel    ALL=(ALL)       NOPASSWD: ALL' > /etc/sudoers.d/nopasswd")

		fmt.Println(prefix, "|", "creating user account")
		goexec.MustRunPrefixed(prefix, "wsl", "-d", distro, "adduser", "-G", "wheel", currentUser)

		fmt.Println(prefix, "|", "setting default user account")
		goexec.MustRunPrefixed(prefix, distroCliTool, "config", "--default-user", currentUser)
	} else {
		fmt.Println(prefix, "|", "user account already exists")
	}

	tempDir, err := ioutil.TempDir("", "bradsDotFiles")
	goerr.Check(err)
	defer os.RemoveAll(tempDir)
	defer fmt.Println(prefix, "|", "deleted", tempDir)
	fmt.Println(prefix, "|", "created", tempDir)

	linuxSetupPath := filepath.Join(tempDir, "setup")
	fmt.Println(prefix, "|", "writing", linuxSetupPath)
	goerr.Check(ioutil.WriteFile(linuxSetupPath, linuxSetupBin, 0777))

	linuxSetupPath = strings.Replace(linuxSetupPath, "C:\\", "/mnt/c/", 1)
	linuxSetupPath = strings.ReplaceAll(linuxSetupPath, "\\", "/")
	goexec.MustRunPrefixed(prefix, "wsl", "-d", distro, "-e", linuxSetupPath)
}

func BoostrapInWSLAsync(distro string) *task.Task {
	return task.New(func() { MustBoostrapInWSL(distro) })
}
