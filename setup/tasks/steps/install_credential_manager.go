// +build windows

package steps

import (
	"fmt"
	"os"

	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/brad-jones/gopwsh"
)

// MustInstallCredentialManager will install https://www.powershellgallery.com/packages/CredentialManager/1.0
// This is used by the "run-at-logon.ps1" script to unlock our SSH/GPG keys without interaction.
func MustInstallCredentialManager() {
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

	// Enforce PowerShell Desktop
	// The CredentialManager module does not work with PS6+
	ps := gopwsh.MustNew(gopwsh.PwshLocation("powershell.exe"))
	defer ps.Exit()

	fmt.Println(prefix, "| installing PowerShellGet")
	ps.MustExecute("Install-Module -Name PowerShellGet -Force")

	fmt.Println(prefix, "| importing PowerShellGet")
	ps.MustExecute("Import-Module PowerShellGet -Force")

	fmt.Println(prefix, "| installing CredentialManager")
	ps.MustExecute("Install-Module CredentialManager -Force")
}

func InstallCredentialManagerAsync() *task.Task {
	return task.New(func() { MustInstallCredentialManager() })
}
