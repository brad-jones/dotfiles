package steps

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/brad-jones/dotfiles/pkg/assets"
	"github.com/brad-jones/dotfiles/pkg/survey"
	"github.com/brad-jones/dotfiles/pkg/tools"
	"github.com/brad-jones/dotfiles/pkg/tools/awsvault"
	"github.com/brad-jones/dotfiles/pkg/tools/chrome"
	"github.com/brad-jones/dotfiles/pkg/tools/dotnet"
	"github.com/brad-jones/dotfiles/pkg/tools/firefox"
	"github.com/brad-jones/dotfiles/pkg/tools/scoop"
	"github.com/brad-jones/dotfiles/pkg/tools/wavebox"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goasync/v2/await"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

// TODO: Install SessionManagerPlugin
// https://s3.amazonaws.com/session-manager-downloads/plugin/latest/windows/SessionManagerPluginSetup.exe
// https://s3.amazonaws.com/session-manager-downloads/plugin/latest/mac/sessionmanager-bundle.zip
// https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_64bit/session-manager-plugin.deb
// https://s3.amazonaws.com/session-manager-downloads/plugin/latest/linux_64bit/session-manager-plugin.rpm

func MustInstallOrUpdate(answers *survey.Answers) {
	await.MustFastAllOrError(
		chrome.InstallAsync(),
		firefox.InstallAsync(),
		wavebox.InstallAsync(),
		awsvault.InstallAsync(),
		dotnet.InstallAsync(tools.DotnetVersions()...),
		task.New(func() {
			if runtime.GOOS == "windows" {
				updateWindows()
				return
			}
			updateLinux()
		}),
	)
}

func InstallOrUpdateAsync(answers *survey.Answers) *task.Task {
	return task.New(func() { MustInstallOrUpdate(answers) })
}

func updateWindows() {
	// Add some additional buckets
	await.MustFastAllOrError(
		scoop.AddBucketAsync("extras", ""),
		scoop.AddBucketAsync("nonportable", ""),
		scoop.AddBucketAsync("java", ""),
		scoop.AddBucketAsync("jetbrains", ""),
		scoop.AddBucketAsync("goreleaser", "https://github.com/goreleaser/scoop-bucket.git"),
		scoop.AddBucketAsync("brad-jones", "https://github.com/brad-jones/scoop-bucket.git"),
	)

	// Now install all the things
	scoop.MustInstallOrUpdatePkgs(map[string]string{
		"7zip":                 "*",
		"adoptopenjdk-hotspot": "*",
		"aws":                  "*",
		"curl":                 "*",
		"dart":                 "*",
		"deno":                 "*",
		"git":                  "*",
		"gitkraken":            "*",
		"go":                   "*",
		"grep":                 "*",
		"jq":                   "*",
		"kotlin":               "*",
		"ktlint":               "*",
		"maven":                "*",
		"nodejs":               "*",
		"nuget":                "*",
		"openssl":              "*",
		"packer":               "*",
		"protobuf":             "*",
		"pwsh":                 "*",
		"python":               "*",
		"ruby":                 "*",
		"sed":                  "*",
		"sonar-scanner":        "*",
		"task":                 "*",
		"terraform":            "*",
		"vlc":                  "*",
		"vscode":               "*",
		"wget":                 "*",
		"windows-terminal":     "*",
	})
}

func updateLinux() {
	prefix := colorchooser.Sprint("update-linux")
	tmpDir, err := ioutil.TempDir("", "bradsDotFiles")
	goerr.Check(err, "failed to create tmpDir")
	defer os.RemoveAll(tmpDir)
	tmpFile := filepath.Join(tmpDir, "dotfiles-updater.sh")
	assets.WriteFile("dotfiles-updater.sh", tmpFile)
	goexec.MustRunPrefixedCmd(prefix, goexec.MustCmd("bash", goexec.Args(tmpFile), goexec.Cwd(utils.HomeDir())))

}
