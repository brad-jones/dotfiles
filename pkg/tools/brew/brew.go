package brew

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/brad-jones/dotfiles/pkg/survey"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

// see: https://docs.brew.sh/Homebrew-on-Linux

var UnSupportedOS = goerr.New("brew not supported by your OS")

func Install(answers *survey.Answers) (err error) {
	defer goerr.Handle(func(e error) { err = e })
	if runtime.GOOS == "windows" {
		goerr.Check(UnSupportedOS)
	}

	// Bail out early if brew already exists
	prefix := colorchooser.Sprint("install-brew")
	if !answers.Reset && utils.CommandExists("brew") {
		fmt.Println(prefix, "|", "already installed, skipping...")
	}

	// Install brew's deps
	utils.RunElevatedNix(prefix, answers.SudoPassword,
		"dnf", "groupinstall", "-y", "Development Tools",
	)
	utils.RunElevatedNix(prefix, answers.SudoPassword, "dnf", "install", "-y",
		"curl",
		"file",
		"git",
		"procps-ng",
		"libxcrypt-compat",
	)

	// Clean up the brew dir if it already exists
	brewDir := filepath.Join(utils.HomeDir(), ".linuxbrew")
	fmt.Println(prefix, "|", "removing", brewDir, "if exists")
	os.RemoveAll(brewDir)

	// Clone the repo
	goexec.MustRunPrefixed(prefix, "git", "clone",
		"https://github.com/Homebrew/brew.git",
		filepath.Join(brewDir, "Homebrew"),
	)

	// Setup the brew symlink as per install docs
	fmt.Println(prefix, "|", "installing brew bin")
	binDir := filepath.Join(brewDir, "bin")
	goerr.Check(os.Mkdir(binDir, 0755), "failing creating", binDir)
	src := filepath.Join(brewDir, "Homebrew/bin/brew")
	dst := filepath.Join(brewDir, "bin/brew")
	goerr.Check(os.Symlink(src, dst), "failed to create symlink", src, dst)

	// Fixes: warning: Insecure world writable dir /usr/bin in PATH, mode 040777
	// Brew is obviously paranoid about these permissions, for good reason I guess.
	// TODO: Why are my permissions so open??? Is it a WSL thing?
	fmt.Println(prefix, "|", "fixing /usr permissions")
	utils.RunElevatedNix(prefix, answers.SudoPassword, "sudo", "chmod", "0755", "/home")
	utils.RunElevatedNix(prefix, answers.SudoPassword, "sudo", "chmod", "0755", "/usr")
	utils.RunElevatedNix(prefix, answers.SudoPassword, "sudo", "chmod", "0755", "/usr/bin")
	utils.RunElevatedNix(prefix, answers.SudoPassword, "sudo", "chmod", "0755", "/usr/sbin")

	// Smoke test the brew command
	goexec.MustRunPrefixed(prefix, Path(), "--version")
	return
}

func MustInstall(answers *survey.Answers) {
	goerr.Check(Install(answers))
}

func InstallAsync(answers *survey.Answers) *task.Task {
	return task.New(func() { MustInstall(answers) })
}

func Path() string {
	return filepath.Join(utils.HomeDir(), ".linuxbrew", "bin", "brew")
}

/*
	rm -rf ~/.linuxbrew;
	sudo dnf install -y curl file git libxcrypt-compat;
	git clone https://github.com/Homebrew/brew ~/.linuxbrew/Homebrew;
	mkdir ~/.linuxbrew/bin;
	ln -s ~/.linuxbrew/Homebrew/bin/brew ~/.linuxbrew/bin;
	eval "$(~/.linuxbrew/bin/brew shellenv)";
	~/.linuxbrew/bin/brew --version;
*/
