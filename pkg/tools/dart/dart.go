package dart

import (
	"path/filepath"
	"runtime"

	"github.com/brad-jones/dotfiles/pkg/tools/brew"
	"github.com/brad-jones/dotfiles/pkg/tools/scoop"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

func Install() (err error) {
	defer goerr.Handle(func(e error) { err = e })
	prefix := colorchooser.Sprint("install-dart")

	if runtime.GOOS == "windows" {
		scoop.MustInstallOrUpdatePkgs(map[string]string{"dart": "*"})
		return
	}

	goexec.MustRunPrefixed(prefix, brew.Path(), "tap", "dart-lang/dart")
	goexec.MustRunPrefixed(prefix, brew.Path(), "install", "dart")

	return
}

func MustInstall() {
	goerr.Check(Install())
}

func InstallAsync() *task.Task {
	return task.New(func() { MustInstall() })
}

func Path() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(scoop.Path(), "apps", "dart", "current", "bin", "dart.exe")
	}
	return filepath.Join(utils.HomeDir(), ".linuxbrew", "bin", "dart")
}

func PubPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(scoop.Path(), "apps", "dart", "current", "bin", "pub.bat")
	}
	return filepath.Join(utils.HomeDir(), ".linuxbrew", "bin", "pub")
}

/*
	rm -rf ~/.dart;
	mkdir ~/.dart;
	tmpFolder="/tmp/$(uuidgen)";
	mkdir -p $tmpFolder;
	function finish {
	rm -rf $tmpFolder;
	}
	curl "https://storage.googleapis.com/dart-archive/channels/stable/release/latest/sdk/dartsdk-linux-x64-release.zip" -o "$tmpFolder/dartsdk.zip";
	unzip "$tmpFolder/dartsdk.zip" -d "$tmpFolder/extracted";
	dartV="$(cat "$tmpFolder/extracted/dart-sdk/version")";
	mv "$tmpFolder/extracted/dart-sdk" "$HOME/.dart/$dartV";
	ln -s "$HOME/.dart/$dartV" "$HOME/.dart/current";
	cd ~/.local/sbin && ~/.dart/current/bin/pub get && cd -;
*/
