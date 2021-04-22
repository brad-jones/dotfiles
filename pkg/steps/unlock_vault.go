package steps

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/brad-jones/dotfiles/pkg/survey"
	"github.com/brad-jones/dotfiles/pkg/tools/gopass"
	"github.com/brad-jones/dotfiles/pkg/tools/gpg"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goasync/v2/await"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/brad-jones/goprefix/v2/pkg/prefixer"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/gosimple/slug"
	"github.com/zalando/go-keyring"
)

const gitUserName = "brad-jones"
const vaultRepoSSH = "git@github.com:brad-jones/vault.git"
const vaultKeyname = "Brad Jones (vault) <brad@bjc.id.au>"
const vaultRepoHTTPS = "https://github.com/brad-jones/vault.git"
const vaultKeyRepoHTTPS = "https://gitlab.com/brad-jones/vault-key.git"

func UnlockVault(answers *survey.Answers) (err error) {
	defer goerr.Handle(func(e error) { err = e })

	// If the tools for the vault exist already but we have no creds then we
	// need to unlock the vault, grab the creds, close the vault, update the
	// vault & finally unlock it again.
	if utils.CommandExists("gpg") &&
		utils.CommandExists("gopass") &&
		len(answers.SudoPassword) == 0 &&
		len(answers.GithubPassword) == 0 &&
		len(answers.GitlabPassword) == 0 {
		gpg.MustStartAgent()
		gpg.MustUnlockKey(vaultKeyname, answers.VaultKeyPassword)
		answers.GithubPassword = gopass.GetSecret("websites/github.com/brad@bjc.id.au")
		answers.GitlabPassword = gopass.GetSecret("websites/gitlab.com/brad@bjc.id.au")
		answers.SudoPassword = gopass.GetSecret(fmt.Sprintf("sudo/%s", slug.Make(utils.GetComputerName())))
		utils.KillProcByName("gpg-agent")
	}

	// Install or update the vault
	await.MustFastAllOrError(
		gopass.InstallAsync(answers.Reset),
		gpg.InstallAsync(answers.SudoPassword),
		downloadVaultAsync(answers.GithubPassword, answers.Reset),
	)

	// Unlock the vault
	gpg.MustStartAgent()
	mustDownloadVaultKey(answers.GitlabPassword, answers.VaultKeyPassword, answers.Reset)
	gpg.MustUnlockKey(vaultKeyname, answers.VaultKeyPassword)

	return
}

func MustUnlockVault(answers *survey.Answers) {
	goerr.Check(UnlockVault(answers))
}

func UnlockVaultAsync(answers *survey.Answers) *task.Task {
	return task.New(func() { MustUnlockVault(answers) })
}

func downloadVaultAsync(repoPassword string, reset bool) *task.Task {
	return task.New(func() {
		prefix := colorchooser.Sprint("download-vault")
		cloneDir := filepath.Join(utils.HomeDir(), ".password-store")

		// Bail out early if the folder exists with a ".git" folder
		if !reset && utils.FolderExists(filepath.Join(cloneDir, ".git")) {
			fmt.Println(prefix, "|", "skipping", cloneDir, "already exists")
			return
		}

		// Delete the .password-store if it exists
		fmt.Println(prefix, "|", "removing", cloneDir, "if it exists")
		goerr.Check(os.RemoveAll(cloneDir), "failed to delete", cloneDir)

		// Setup the prefixer for the following clone op
		r, w, err := os.Pipe()
		goerr.Check(err, "failed to create os.Pipe")
		go prefixer.New(prefix + " | ").ReadFrom(r).WriteTo(os.Stdout)

		// Clone the git repo into the .password-store dir
		fmt.Println(prefix, "|", "cloning", vaultRepoHTTPS, "into", cloneDir)
		repo, err := git.PlainClone(cloneDir, false, &git.CloneOptions{
			URL:      vaultRepoHTTPS,
			Progress: w,
			Auth: &http.BasicAuth{
				Username: gitUserName,
				Password: repoPassword,
			},
		})
		goerr.Check(err, "failed to clone", vaultRepoHTTPS)

		// All future git ops on this repo should use SSH instead,
		// the appropriate SSH keys will be installed later on by this app.
		fmt.Println(prefix, "|", "reconfiguring the origin to use SSH in the future")
		c, err := repo.Config()
		goerr.Check(err, "failed to get repo config")
		c.Remotes["origin"].URLs = []string{vaultRepoSSH}
		goerr.Check(repo.SetConfig(c), "failed to set repo config")
	})
}

func mustDownloadVaultKey(repoPassword, keyPassword string, reset bool) {
	prefix := colorchooser.Sprint("download-vault-key")

	// Bail out if the key already exists
	if !reset && gpg.PrivateKeyExists(vaultKeyname) {
		fmt.Println(prefix, "skipping, key already exists")
		return
	}

	// Create a tmp dir to clone into
	cloneDir, err := ioutil.TempDir("", "bradsDotFiles")
	goerr.Check(err, "failed to create temp dir")
	defer os.RemoveAll(cloneDir)
	defer fmt.Println(prefix, "|", "deleting", cloneDir)

	// Setup the prefixer for the following clone op
	r, w, err := os.Pipe()
	goerr.Check(err, "failed to create os.Pipe")
	go prefixer.New(prefix + " | ").ReadFrom(r).WriteTo(os.Stdout)

	// Clone the git repo into the temp dir
	fmt.Println(prefix, "|", "cloning", vaultKeyRepoHTTPS, "into", cloneDir)
	_, err = git.PlainClone(cloneDir, false, &git.CloneOptions{
		URL:      vaultKeyRepoHTTPS,
		Progress: w,
		Auth: &http.BasicAuth{
			Username: gitUserName,
			Password: repoPassword,
		},
	})
	goerr.Check(err, "failed to clone", vaultKeyRepoHTTPS)

	// Import the key into the gpg keychain
	gpg.MustImportKey(filepath.Join(cloneDir, "private.pem"),
		vaultKeyname,
		keyPassword,
		false,
	)

	// On windows we will store the passphrase into the wincred store
	// On other systems, like a Fedora desktop, we rely on the built-in SSH/GPG
	// passphrase caching and automatic unlocking.
	if runtime.GOOS == "windows" {
		fmt.Println(prefix, "|", "storing the vault key passphrase into local wincred store")
		goerr.Check(keyring.Set("passphrase", "vault", keyPassword),
			"failed to store vault passphrase into native OS keychain",
		)
	}
}
