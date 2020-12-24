package steps

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/cavaliercoder/grab"
	"github.com/google/go-github/github"
	"github.com/mholt/archiver"
	"github.com/phayes/permbits"
)

// InstallGithubPkg will download, extract and place into your home dir a binary from a github release
func InstallGithubPkg(owner, repo, tag, exeName string) {
	g := github.NewClient(nil)
	prefix := colorchooser.Sprint("install-"+repo) + " |"

	fmt.Println(prefix, fmt.Sprintf("getting release assets from github.com/%s/%s/releases/tag/%s", owner, repo, tag))
	r, _, err := g.Repositories.GetReleaseByTag(context.Background(), owner, repo, tag)
	goerr.Check(err, owner, repo, tag)

	downloadURL := ""
	for _, v := range r.Assets {
		name := strings.ToLower(v.GetName())
		if isArchive(name) && strings.Contains(name, runtime.GOOS) && strings.Contains(name, runtime.GOARCH) {
			downloadURL = v.GetBrowserDownloadURL()
			break
		}
	}
	if downloadURL == "" {
		goerr.Check(Pkg404Error, runtime.GOOS, owner, repo, tag)
	}

	tmpDir, err := ioutil.TempDir("", repo)
	goerr.Check(err, repo)
	defer os.RemoveAll(tmpDir)
	defer func() { fmt.Println(prefix, "deleting", tmpDir) }()

	fmt.Println(prefix, "downloading", downloadURL, "into", tmpDir)
	resp, err := grab.Get(filepath.Join(tmpDir, "."), downloadURL)
	goerr.Check(err, downloadURL, tmpDir)

	fmt.Println(prefix, "extracting", resp.Filename)
	extracted := filepath.Join(tmpDir, "extracted")
	goerr.Check(archiver.Unarchive(resp.Filename, extracted), resp.Filename, extracted)

	home, err := os.UserHomeDir()
	goerr.Check(err)

	src := filepath.Join(extracted, exeName)
	dst := filepath.Join(home, ".local", "bin", exeName)
	if runtime.GOOS == "windows" {
		dst = dst + ".exe"
		src = src + ".exe"
	}

	fmt.Println(prefix, "moving", src, "to", dst)
	goerr.Check(os.MkdirAll(filepath.Dir(dst), 0644), dst)
	goerr.Check(os.Rename(src, dst), src, dst)

	if runtime.GOOS != "windows" {
		fmt.Println(prefix, "setting execute bit for", dst)

		permissions, err := permbits.Stat(dst)
		goerr.Check(err)
		permissions.SetUserExecute(true)
		permissions.SetGroupExecute(true)
		permissions.SetOtherExecute(true)
		goerr.Check(permbits.Chmod(dst, permissions), dst, permissions.String())
	}
}

func isArchive(filename string) bool {
	return strings.HasSuffix(filename, ".zip") ||
		strings.HasSuffix(filename, ".tar.gz")
}
