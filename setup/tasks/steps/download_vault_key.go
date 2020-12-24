package steps

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/avast/retry-go"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

// DownloadVaultKey will fetch the key from the gitlab repo and import it into the keychain
func DownloadVaultKey(repoPassword, keyPassword string) {
	prefix := colorchooser.Sprint("download-vault-key")

	cloneDir, err := ioutil.TempDir("", "bradsDotFiles")
	goerr.Check(err)
	defer os.RemoveAll(cloneDir)

	// TODO: Don't leak password here - use GIT_ASKPASS?
	cloneURI := "https://brad-jones:" + url.QueryEscape(repoPassword) + "@gitlab.com/brad-jones/vault-key.git"
	goexec.MustRunPrefixed(prefix, "git", "clone", cloneURI, cloneDir)

	// Import the gpg key into the keychain
	goexec.MustRunPrefixed(prefix, "gpg", "--import", filepath.Join(cloneDir, "private.pem"))

	// Trust the key
	trustGpgKey(prefix, "Brad Jones (vault) <brad@bjc.id.au>")

	// Add the key to the agent
	if runtime.GOOS == "windows" {
		fmt.Println(prefix, "starting gpg agent...")
		goerr.Check(retry.Do(func() error {
			return goexec.RunPrefixed(prefix, "gpg-connect-agent", "/bye")
		}))

		fmt.Println(prefix, "adding preset for 83D182028C7F2DF102F09E61FF308BBB10F539D8")
		goexec.MustRunPrefixed(prefix, "gpg-preset-passphrase",
			"--passphrase", keyPassword,
			"--preset", "83D182028C7F2DF102F09E61FF308BBB10F539D8", // keygrip
		)

		fmt.Println(prefix, "adding preset for F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB")
		goexec.MustRunPrefixed(prefix, "gpg-preset-passphrase",
			"--passphrase", keyPassword,
			"--preset", "F217E464BDDC0DF42C0E4B5F740FD611F4E35ADB", // keygrip
		)
	}
}
