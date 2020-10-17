package tasks

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/brad-jones/goerr/v2"
	"github.com/cavaliercoder/grab"
	"github.com/google/go-github/v32/github"
	"github.com/mholt/archiver/v3"
	"github.com/phayes/permbits"
)

// Bootstrap will run when someone executes this program directly without
// any additional input and it's job is to do all the things required to
// setup chezmoi.
func Bootstrap() error {
	if err := installChezmoi(); err != nil {
		return goerr.Wrap(err)
	}

	if err := installGoPass(); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}

func installChezmoi() error {
	g := github.NewClient(nil)

	fmt.Println("getting latest release from github.com/twpayne/chezmoi")
	r, _, err := g.Repositories.GetLatestRelease(
		context.Background(), "twpayne", "chezmoi",
	)
	if err != nil {
		return goerr.Wrap(err)
	}

	downloadURL := ""
	for _, v := range r.Assets {
		if runtime.GOOS == "windows" {
			if strings.Contains(v.GetName(), "windows_"+runtime.GOARCH+".zip") {
				downloadURL = v.GetBrowserDownloadURL()
				break
			}
		} else {
			if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
				if strings.Contains(v.GetName(), runtime.GOOS+"_"+runtime.GOARCH+".tar.gz") {
					downloadURL = v.GetBrowserDownloadURL()
					break
				}
			}
		}
	}
	if downloadURL == "" {
		return goerr.Wrap("could not locate download url for chezmoi")
	}

	tmpDir, err := ioutil.TempDir("", "ChezmoiSetup")
	if err != nil {
		return goerr.Wrap(err)
	}
	defer os.RemoveAll(tmpDir)

	fmt.Println("downloading", downloadURL, "into", tmpDir)
	resp, err := grab.Get(filepath.Join(tmpDir, "."), downloadURL)
	if err != nil {
		return goerr.Wrap(err)
	}

	fmt.Println("extracting", resp.Filename)
	extracted := filepath.Join(tmpDir, "extracted")
	if err := archiver.Unarchive(resp.Filename, extracted); err != nil {
		return goerr.Wrap(err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return goerr.Wrap(err)
	}

	src := filepath.Join(extracted, "chezmoi")
	dst := filepath.Join(home, ".local", "bin", "chezmoi")
	if runtime.GOOS == "windows" {
		dst = dst + ".exe"
		src = src + ".exe"
	}

	fmt.Println("moving", src, "to", dst)

	if err := os.MkdirAll(filepath.Dir(dst), 0644); err != nil {
		return goerr.Wrap(err)
	}

	if err := os.Rename(src, dst); err != nil {
		return goerr.Wrap(err)
	}

	if runtime.GOOS != "windows" {
		fmt.Println("setting execute bit for", dst)

		permissions, err := permbits.Stat(dst)
		if err != nil {
			return goerr.Wrap(err)
		}
		permissions.SetUserExecute(true)
		permissions.SetGroupExecute(true)
		permissions.SetOtherExecute(true)
		if err := permbits.Chmod(dst, permissions); err != nil {
			return goerr.Wrap(err)
		}
	}

	fmt.Println("chezmoi is installed")

	return nil
}

func installGoPass() error {
	g := github.NewClient(nil)

	fmt.Println("getting latest release from github.com/gopasspw/gopass")
	r, _, err := g.Repositories.GetReleaseByTag(
		context.Background(), "gopasspw", "gopass", "v1.9.2",
	)
	if err != nil {
		return goerr.Wrap(err)
	}

	downloadURL := ""
	for _, v := range r.Assets {
		if runtime.GOOS == "windows" {
			if strings.Contains(v.GetName(), "windows-"+runtime.GOARCH+".zip") {
				downloadURL = v.GetBrowserDownloadURL()
				break
			}
		} else {
			if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
				if strings.Contains(v.GetName(), runtime.GOOS+"-"+runtime.GOARCH+".tar.gz") {
					downloadURL = v.GetBrowserDownloadURL()
					break
				}
			}
		}
	}
	if downloadURL == "" {
		return goerr.Wrap("could not locate download url for gopass")
	}

	tmpDir, err := ioutil.TempDir("", "GoPassSetup")
	if err != nil {
		return goerr.Wrap(err)
	}
	defer os.RemoveAll(tmpDir)

	fmt.Println("downloading", downloadURL, "into", tmpDir)
	resp, err := grab.Get(filepath.Join(tmpDir, "."), downloadURL)
	if err != nil {
		return goerr.Wrap(err)
	}

	fmt.Println("extracting", resp.Filename)
	extracted := filepath.Join(tmpDir, "extracted")
	if err := archiver.Unarchive(resp.Filename, extracted); err != nil {
		return goerr.Wrap(err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return goerr.Wrap(err)
	}

	src := filepath.Join(extracted, "gopass")
	dst := filepath.Join(home, ".local", "bin", "gopass")
	if runtime.GOOS == "windows" {
		dst = dst + ".exe"
		src = src + ".exe"
	}

	fmt.Println("moving", src, "to", dst)

	if err := os.MkdirAll(filepath.Dir(dst), 0644); err != nil {
		return goerr.Wrap(err)
	}

	if err := os.Rename(src, dst); err != nil {
		return goerr.Wrap(err)
	}

	if runtime.GOOS != "windows" {
		fmt.Println("setting execute bit for", dst)

		permissions, err := permbits.Stat(dst)
		if err != nil {
			return goerr.Wrap(err)
		}
		permissions.SetUserExecute(true)
		permissions.SetGroupExecute(true)
		permissions.SetOtherExecute(true)
		if err := permbits.Chmod(dst, permissions); err != nil {
			return goerr.Wrap(err)
		}
	}

	fmt.Println("gopass is installed")

	return nil
}
