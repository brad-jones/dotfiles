package steps

import (
	"github.com/brad-jones/goexec/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

// InstallPSReadLine will install https://github.com/PowerShell/PSReadLine
//
// Highly likely to fail due to all this shenanigans:
// - https://github.com/PowerShell/PSReadLine#upgrading
// - https://github.com/PowerShell/PSReadLine/issues/1370
//
// Probs just going to leave this disabled for now as there is no easy way to solve this
func InstallPSReadLine() {
	prefix := colorchooser.Sprint("install-ps-readline")
	goexec.MustRunPrefixed(prefix, "pwsh", "-NoProfile", "-Command",
		"Install-Module PSReadLine -Force -SkipPublisherCheck -AllowPrerelease",
	)
}
