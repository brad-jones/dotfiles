package ghpkg

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/cavaliercoder/grab"
	"github.com/google/go-github/v35/github"
	"github.com/mholt/archiver"
	"github.com/phayes/permbits"
)

/*
	Thought bubble:
	This is almost worthy of being a totally standalone project / tool
*/

// Pkg404Error is returned when we could not locate a suitable URL to download
var Pkg404Error = goerr.New("could not locate download url")

// PkgInvalidHash is returned when the hashes of downloaded binary do not match with the supplied hash
var PkgInvalidHash = goerr.New("downloaded artifact does not match supplied hash")

// We are using the Functional Options pattern.
// see: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type Config struct {
	owner      string
	repo       string
	tag        string
	exeName    string
	dstExeName string
	sha256Hash string
	pkgPattern string
	reset      bool
	naked      bool
}

// In addition to a git tag, you can set a SHA256 hash to make sure
// what was downloaded really is cryptographically what you expected.
func Sha256Hash(value string) func(*Config) error {
	return func(c *Config) error {
		c.sha256Hash = value
		return nil
	}
}

// By default we apply some simple logic based on the values of GOOS & GOARCH
// to determin which Github Asset we should download. For most packages this
// works ok but sometimes we need to be more explicit.
func PkgPattern(value string) func(*Config) error {
	return func(c *Config) error {
		c.pkgPattern = value
		return nil
	}
}

// The name of the actual binary inside the downloaded archive.
// If not set we assume it's the same as the repo name.
func ExeName(value string) func(*Config) error {
	return func(c *Config) error {
		c.exeName = value
		return nil
	}
}

// Sometimes the name of binary inside the archive is different than what you
// want it to be in the final installed location. This allows you to effectively
// rename the binary. If not set the value of ExeName is used.
func DstExeName(value string) func(*Config) error {
	return func(c *Config) error {
		c.dstExeName = value
		return nil
	}
}

// If set to true then, the binary will be downloaded and re-installed even if
// the hashes match and everything else looks good.
func Reset(value bool) func(*Config) error {
	return func(c *Config) error {
		c.reset = value
		return nil
	}
}

func Naked(value bool) func(*Config) error {
	return func(c *Config) error {
		c.naked = value
		return nil
	}
}

// InstallPkg downloads released archives from Github Releases, of a given repo,
// extracts the contents and then copies a single self contained binary out of
// the package and into ~/.local/bin.
func InstallPkg(owner, repo, tag string, decorators ...func(*Config) error) (err error) {
	defer goerr.Handle(func(e error) { err = e })

	// Set defaults and read in our config
	c := &Config{
		owner:   owner,
		repo:    repo,
		tag:     tag,
		exeName: repo,
	}
	for _, d := range decorators {
		goerr.Check(d(c), "decorator failed")
	}

	g := github.NewClient(nil)
	prefix := colorchooser.Sprint("install-" + repo)

	// Figure out what the destination path is for the binary
	dst := filepath.Join(utils.HomeDir(), ".local", "bin", c.exeName)
	if len(c.dstExeName) > 0 {
		dst = filepath.Join(utils.HomeDir(), ".local", "bin", c.dstExeName)
	}
	if runtime.GOOS == "windows" && !strings.HasSuffix(dst, ".exe") {
		dst = dst + ".exe"
	}

	// If the destination already exists and is the same hash then bail out early
	if !c.reset && len(c.sha256Hash) > 0 && utils.FileExists(dst) {
		if utils.Sha256HashFile(dst) == c.sha256Hash {
			fmt.Println(prefix, "|", "skipping, already installed")
			return
		}
	}

	// Get a list of release assets from github
	fmt.Println(prefix, "|", fmt.Sprintf("getting release assets from github.com/%s/%s/releases/tag/%s", owner, repo, tag))
	r, _, err := g.Repositories.GetReleaseByTag(context.Background(), owner, repo, tag)
	goerr.Check(err, "GetReleaseByTag failed", owner, repo, tag)

	// Loop through those assets looking for a valid package to download
	downloadURL := ""
	for _, v := range r.Assets {
		name := strings.ToLower(v.GetName())
		if isDownloadable(name, c.pkgPattern) {
			downloadURL = v.GetBrowserDownloadURL()
			break
		}
	}
	if downloadURL == "" {
		goerr.Check(Pkg404Error, runtime.GOOS, owner, repo, tag)
	}

	// Create a temp dir to download the package into
	tmpDir, err := ioutil.TempDir("", repo)
	goerr.Check(err, "failed to create temp dir", repo)
	defer os.RemoveAll(tmpDir)
	defer fmt.Println(prefix, "|", "deleting", tmpDir)

	// Download the archive into the temp dir
	fmt.Println(prefix, "|", "downloading", downloadURL, "into", tmpDir)
	resp, err := grab.Get(filepath.Join(tmpDir, "."), downloadURL)
	goerr.Check(err, "failed to download", downloadURL, tmpDir)

	var src string
	if c.naked {
		src = resp.Filename
	} else {
		// Extract the archive
		extracted := filepath.Join(tmpDir, "extracted")
		fmt.Println(prefix, "|", "extracting", resp.Filename)
		goerr.Check(archiver.Unarchive(resp.Filename, extracted),
			"failed to extract archive", resp.Filename, extracted,
		)

		// Work out what the binary name is inside the archive
		src = filepath.Join(extracted, c.exeName)
		if runtime.GOOS == "windows" && !strings.HasSuffix(src, ".exe") {
			src = src + ".exe"
		}
	}

	// Check to make sure it matches our hash
	if len(c.sha256Hash) > 0 {
		if utils.Sha256HashFile(src) != c.sha256Hash {
			goerr.Check(PkgInvalidHash, owner, repo, tag, c.sha256Hash)
		}
	}

	// Install the binary in the final location
	fmt.Println(prefix, "|", "moving", src, "to", dst)
	goerr.Check(os.MkdirAll(filepath.Dir(dst), 0755), dst)
	goerr.Check(os.Rename(src, dst), src, dst)

	// On *nix make sure it is executable
	if runtime.GOOS != "windows" {
		fmt.Println(prefix, "|", "setting execute bit for", dst)

		permissions, err := permbits.Stat(dst)
		goerr.Check(err, "failed to stat permissions")
		permissions.SetUserExecute(true)
		permissions.SetGroupExecute(true)
		permissions.SetOtherExecute(true)
		goerr.Check(permbits.Chmod(dst, permissions),
			"failed to set execute bit",
			dst, permissions.String(),
		)
	}

	return
}

// MustInstallPkg does the same thing InstallPkg but panics instead of returning an error
func MustInstallPkg(owner, repo, tag string, decorators ...func(*Config) error) {
	goerr.Check(InstallPkg(owner, repo, tag, decorators...))
}

// InstallPkgAsync does the same thing as InstallPkg but asynchronously.
func InstallPkgAsync(owner, repo, tag string, decorators ...func(*Config) error) *task.Task {
	return task.New(func() { MustInstallPkg(owner, repo, tag, decorators...) })
}

func isArchive(filename string) bool {
	return strings.HasSuffix(filename, ".zip") ||
		strings.HasSuffix(filename, ".tar.gz")
}

func isDownloadable(name, pkgPattern string) bool {
	if len(pkgPattern) > 0 {
		match, err := regexp.MatchString(pkgPattern, name)
		goerr.Check(err, pkgPattern, "did not match", name)
		return match
	}

	if !isArchive(name) {
		return false
	}

	if !strings.Contains(name, runtime.GOARCH) {
		return false
	}

	if !strings.Contains(name, runtime.GOOS) {
		return false
	}

	return true
}
