// +build windows

package steps

import (
	"fmt"

	"github.com/brad-jones/dotfiles/setup/tasks/utils"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/brad-jones/gopwsh"
)

// MustDisableUACForWindows will disable UAC.
//
// I liken UAC to SELinux, I have gotten by without both technologies and
// I am yet to have a security incident. Perhaps thats naive but oh well...
//
// This does not disable UAC completely, we do not set `EnableLUA` to 0.
// As this would break vital built-in components of Windows, like the Firewall.
//
// We just stop the prompt from displaying for Admin users with
// `ConsentPromptBehaviorAdmin` set to 0.
func MustDisableUACForWindows() {
	prefix := colorchooser.Sprint("disable-uac")
	fmt.Println(prefix, "| setting ConsentPromptBehaviorAdmin to 0")

	ps := gopwsh.MustNew(gopwsh.Elevated(utils.SudoBin()))
	defer ps.Exit()
	ps.MustExecute(`Set-ItemProperty -Path REGISTRY::HKEY_LOCAL_MACHINE\Software\Microsoft\Windows\CurrentVersion\Policies\System -Name ConsentPromptBehaviorAdmin -Value 0`)

	fmt.Println(prefix, "| ConsentPromptBehaviorAdmin is now 0")
}

func DisableUACForWindowsAsync() *task.Task {
	return task.New(func() { MustDisableUACForWindows() })
}
