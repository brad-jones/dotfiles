package wsl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/brad-jones/dotfiles/pkg/assets"
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
func InstallFedora(version, hash, sudoPassword string, makeDefault, reset bool) (name string, err error) {
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

	// There are a collection of tools that you kinda just expect to exist.
	// But because we are starting from essentially the fedora docker image,
	// many of these tools are not installed so lets install them now.
	goexec.MustRunPrefixed(prefix, "wsl", "-d", name,
		"dnf", "groupinstall", "-y", "Development Tools",
	)
	goexec.MustRunPrefixed(prefix, "wsl", "-d", name, "dnf", "install", "-y",
		"cracklib-dicts",
		"curl",
		"findutils",
		"iputils",
		"passwd",
		"procps-ng",
		"sudo",
		"tar",
		"unzip",
		"wget",
		"which",
	)

	// Next is another container issue, users don't have permissions to ping.
	// Thats super annoying when debugging network things so lets make sure
	// we can easily ping stuff.
	//
	// see: https://stackoverflow.com/questions/28553923
	fmt.Println(prefix, "|", "fix ping permissions")
	goexec.MustRunPrefixed(prefix, "wsl", "-d", name,
		"chmod", "4755", "/bin/ping",
	)

	// There is a really annoying issue with WSL2 and the go dns resolver.
	// I seem to be able to fix this by setting a custom set of resolvers.
	//
	// see: https://github.com/golang/go/issues/44135
	// also: https://superuser.com/questions/1533291
	fmt.Println(prefix, "|", "fix dns resolution")
	wslConf := strings.ReplaceAll(string(assets.ReadFile("etc/wsl.conf")), "\n", `\n`)
	goexec.MustRunPrefixed(prefix, "wsl", "-d", name,
		"sh", "-c", fmt.Sprintf(`echo -e "%s" > /etc/wsl.conf`, wslConf),
	)
	fmt.Println(prefix, "|", "written /etc/wsl.conf")
	resolvConf := strings.ReplaceAll(string(assets.ReadFile("etc/resolv.conf")), "\n", `\n`)
	goexec.MustRunPrefixed(prefix, "wsl", "-d", name,
		"sh", "-c", fmt.Sprintf(`echo -e "%s" > /etc/resolv.conf`, resolvConf),
	)
	fmt.Println(prefix, "|", "written /etc/resolv.conf")
	fmt.Println(prefix, "|", "shutting down wsl")
	goexec.MustRunPrefixed(prefix, "wsl", "--shutdown")

	// Here we create a matching user account inside the WSL instance
	currentUser := os.Getenv("USERNAME")
	currentWSLUser := strings.TrimSpace(goexec.MustRunBuffered("wsl", "-d", name, "echo", "$USER").StdOut)
	if currentUser != currentWSLUser {
		fmt.Println(prefix, "|", "creating user account")
		goexec.MustRunPrefixed(prefix, "wsl", "-d", name, "adduser", "-G", "wheel", currentUser)

		fmt.Println(prefix, "|", "setting sudo password for user account")
		goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd("wsl",
			goexec.SetIn(strings.NewReader(sudoPassword)),
			goexec.Args("-d", name, "passwd", currentUser, "--stdin"),
		))

		fmt.Println(prefix, "|", "setting default user account")
		distroCliTool := filepath.Join(utils.HomeDir(), ".wsl", name, name+".exe")
		goexec.MustRunPrefixed(prefix, distroCliTool, "config", "--default-user", currentUser)
	} else {
		fmt.Println(prefix, "|", "user account already exists")
	}

	return
}

func MustInstallFedora(version, hash, sudoPassword string, makeDefault, reset bool) string {
	name, err := InstallFedora(version, hash, sudoPassword, makeDefault, reset)
	goerr.Check(err)
	return name
}

func InstallFedoraAsync(version, hash, sudoPassword string, makeDefault, reset bool) *task.Task {
	return task.New(func() { MustInstallFedora(version, hash, sudoPassword, makeDefault, reset) })
}
