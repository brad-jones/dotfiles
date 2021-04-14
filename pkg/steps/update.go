package steps

import (
	"runtime"

	"github.com/brad-jones/dotfiles/pkg/tools"
	"github.com/brad-jones/dotfiles/pkg/tools/chrome"
	"github.com/brad-jones/dotfiles/pkg/tools/dotnet"
	"github.com/brad-jones/dotfiles/pkg/tools/firefox"
	"github.com/brad-jones/dotfiles/pkg/tools/scoop"
	"github.com/brad-jones/dotfiles/pkg/tools/wavebox"
	"github.com/brad-jones/goasync/v2/await"
	"github.com/brad-jones/goasync/v2/task"
)

func MustUpdate() {
	await.MustFastAllOrError(
		chrome.InstallAsync(),
		firefox.InstallAsync(),
		wavebox.InstallAsync(),
		dotnet.InstallAsync("latest", "3.1.407", "2.1.814"),
		task.New(func() {
			if runtime.GOOS == "windows" {
				updateWindows()
				return
			}
			updateLinux()
		}),
	)
}

func UpdateAsync() *task.Task {
	return task.New(func() { MustUpdate() })
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
		"aws-vault":            tools.GetVersion("aws-vault").No,
		"aws":                  "*",
		"curl":                 "*",
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

func updateLinux() {}
