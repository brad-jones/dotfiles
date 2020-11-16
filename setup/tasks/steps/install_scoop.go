package steps

import (
	"fmt"

	"github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

// InstallScoop will install the https://scoop.sh/ package manager for Windows
//
// TODO: change installation dir to ~/.scoop
func InstallScoop() {
	prefix := colorchooser.Sprint("install-scoop")

	ps, err := powershell.New(&backend.Local{})
	goerr.Check(err)
	defer ps.Exit()

	if _, _, err := ps.Execute("Get-Command scoop"); err == nil {
		goexec.MustRunPrefixed(prefix, "powershell",
			"-Command", "scoop update",
		)
		return
	}

	fmt.Println(prefix, "| setting execution policy to RemoteSigned")
	_, _, err = ps.Execute("Set-ExecutionPolicy RemoteSigned -scope CurrentUser")
	goerr.Check(err)

	goexec.MustRunPrefixed(prefix, "powershell",
		"-Command", "iwr -useb get.scoop.sh | iex",
	)
}
