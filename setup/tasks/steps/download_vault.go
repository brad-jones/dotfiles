package steps

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

// DownloadVault will clone the gopass vault repo into ~/.password-store
func DownloadVault(repoPassword string) {
	prefix := colorchooser.Sprint("download-vault")

	homeDir, err := os.UserHomeDir()
	goerr.Check(err)
	cloneDir := filepath.Join(homeDir, ".password-store")

	fmt.Println(prefix, "removing", cloneDir, "if it exists")
	goerr.Check(os.RemoveAll(cloneDir), cloneDir)

	// TODO: Don't leak password here - use GIT_ASKPASS?
	cloneURI := "https://brad-jones:" + url.QueryEscape(repoPassword) + "@github.com/brad-jones/vault.git"
	goexec.MustRunPrefixed(prefix, "git", "clone", cloneURI, cloneDir)
	goexec.MustRunPrefixed(prefix, "git",
		"--git-dir", filepath.Join(cloneDir, ".git"),
		"remote", "set-url", "origin", "git@github.com:brad-jones/vault.git",
	)
}
