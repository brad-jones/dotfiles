package steps

import (
	"fmt"
	"os"

	"github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

// InstallCredentialManager will install https://www.powershellgallery.com/packages/CredentialManager/1.0
// This is used by the "run-at-logon.ps1" script to unlock our SSH/GPG keys without interaction.
func InstallCredentialManager() {
	prefix := colorchooser.Sprint("install-credential-manager")

	homeDir, err := os.UserHomeDir()
	goerr.Check(err)

	oldValue := os.Getenv("PSModulePath")
	defer os.Setenv("PSModulePath", oldValue)

	goerr.Check(os.Setenv("PSModulePath",
		homeDir+"\\Documents\\WindowsPowerShell\\Modules;"+
			"C:\\Program Files\\WindowsPowerShell\\Modules;"+
			"C:\\WINDOWS\\system32\\WindowsPowerShell\\v1.0\\Modules",
	))

	ps, err := powershell.New(&backend.Local{})
	goerr.Check(err)
	defer ps.Exit()

	fmt.Println(prefix, "| installing PowerShellGet")
	_, _, err = ps.Execute("Install-Module -Name PowerShellGet -Force")
	goerr.Check(err)

	fmt.Println(prefix, "| importing PowerShellGet")
	_, _, err = ps.Execute("Import-Module PowerShellGet -Force")
	goerr.Check(err)

	fmt.Println(prefix, "| installing CredentialManager")
	_, _, err = ps.Execute("Install-Module CredentialManager -Force")
	goerr.Check(err)
}
