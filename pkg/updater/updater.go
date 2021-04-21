package updater

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/brad-jones/dotfiles/pkg/survey"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/dotfiles/pkg/utils/downloader"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-github/v35/github"
	"github.com/tidwall/gjson"
)

const owner = "brad-jones"
const repo = "dotfiles"
const signer = "0x33dc7b56c2be6175e1ad17e31f003f55943fa4ce"

var UnTrusted = goerr.New("the downloaded binary did not authenticate agains codenotary.io")

func Update(currentVersion string, answers *survey.Answers) (err error) {
	defer goerr.Handle(func(e error) { err = e })
	prefix := colorchooser.Sprint("self-update")

	g := github.NewClient(nil)

	// If we don't have a specific version to update to,
	// lets get the latest version from Github
	if answers.UpdateToVersion == "" {
		fmt.Println(prefix, "|", "getting latest release")
		r, _, err := g.Repositories.GetLatestRelease(context.Background(), owner, repo)
		goerr.Check(err, "GetLatestRelease failed")
		answers.UpdateToVersion = r.GetTagName()
		fmt.Println(prefix, "|", "latest release is", answers.UpdateToVersion)
	}

	// Parse the current & new version numbers
	curV, err := semver.Parse(stripV(currentVersion))
	goerr.Check(err, "currentVersion is not a valid semver", currentVersion)
	newV, err := semver.Parse(stripV(answers.UpdateToVersion))
	goerr.Check(err, "answers.UpdateToVersion is not a valid semver", answers.UpdateToVersion)

	// Compare them, if they are the same then we have nothing to do.
	if curV.Equals(newV) {
		fmt.Println(prefix, "|", "skipping, already latest version")
		return
	}

	// Download the new binary to a temp location
	fmt.Println(prefix, "|", "updating from", curV, "to", newV)
	tmpDir, err := ioutil.TempDir("", "bradsDotFiles")
	goerr.Check(err, "failed to create temp dir", repo)
	defer os.RemoveAll(tmpDir)
	defer fmt.Println(prefix, "|", "deleting", tmpDir)
	downloadURL := fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/v%s/dotfiles_%s_%s%s",
		owner, repo, newV, runtime.GOOS, runtime.GOARCH, exeExtention(),
	)
	newBin := downloader.MustDownloadWithProgress(prefix, downloadURL, filepath.Join(tmpDir, "."))

	// Authenticate the new bin against CodeNotary.io
	fmt.Println(prefix, "|", "authenticating against CodeNotary.io")
	resp, err := resty.New().R().
		SetPathParam("hash", utils.Sha256HashFile(newBin)).
		SetQueryParam("signers", signer).
		Get("https://api.codenotary.io/authenticate/{hash}")
	goerr.Check(err, "failed to make vcn request")
	j := gjson.ParseBytes(resp.Body())
	statusCode := j.Get("verification.status").Int()
	if statusCode != 0 {
		goerr.Check(UnTrusted)
	}
	if j.Get("name").String() != downloadURL {
		goerr.Check(UnTrusted, "Downloaded URL does not match Notarized URL from codenotary.io")
	}
	fmt.Println(prefix, "|", "TRUSTED")

	// Now we execute the new version
	fmt.Println(prefix, "|", "running new version")
	if err := goexec.Run(newBin); err != nil {
		// Something went wrong, lets essentially do nothing, the binary will
		// be deleted when the temp dir is cleaned up by our deferred statements
		// above. And then when this returns the rest of this versions process
		// will run, hopefully restoring the system into a working state.
		fmt.Println(prefix, "|", "the new version failed, rolling back...")
		return nil
	}

	// Now we can finally replace and update ourselves.
	fmt.Println(prefix, "|", "new version has run successfully, replacing ourselves")
	goerr.Check(replaceSelf(newBin), "failed to replace ourselves, perform update manually...")
	fmt.Println(prefix, "|", "new version has been installed")

	return
}

func MustUpdate(currentVersion string, answers *survey.Answers) {
	goerr.Check(Update(currentVersion, answers))
}

func stripV(in string) string {
	return strings.TrimPrefix(in, "v")
}

func exeExtention() string {
	exe := ""
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}
	return exe
}

// this was essentially copied from https://github.com/inconshreveable/go-update
func replaceSelf(newCandidate string) (err error) {
	defer goerr.Handle(func(e error) { err = e })

	currentPath, err := os.Executable()
	goerr.Check(err, "failed to get path to current executeable")
	oldPath := currentPath + ".old"

	// delete any existing old exec file - this is necessary on Windows for two reasons:
	// 1. after a successful update, Windows can't remove the .old file because the process is still running
	// 2. windows rename operations fail if the destination file already exists
	os.Remove(oldPath)

	// move the existing executable to a new file in the same directory
	goerr.Check(os.Rename(currentPath, oldPath),
		"failed to move current to old", currentPath, oldPath,
	)

	// move the new executable in to become the new program
	if err := os.Rename(newCandidate, currentPath); err != nil {
		// move unsuccessful
		//
		// The filesystem is now in a bad state. We have successfully moved the
		// existing binary to a new location, but we couldn't move the new binary
		// to take its place. That means there is no file where the current executable
		// binary used to be!
		//
		// Try to rollback by restoring the old binary to its original path.
		goerr.Check(os.Rename(oldPath, currentPath),
			"failed to move old to current", oldPath, currentPath,
		)

		return goerr.Wrap(err, "failed to move new to current", newCandidate, currentPath)
	}

	// move successful, remove the old binary if needed
	errRemove := os.Remove(oldPath)

	// windows has trouble with removing old binaries, so hide it instead
	if errRemove != nil {
		// if the hide fails it's not the end of the world, the rest of the update has worked
		hideFile(oldPath)
	}

	return
}
