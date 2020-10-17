package tasks

import (
	"fmt"
)

// ChezmoiApply will run when chezmoi it's self executes
// this program via the `run_setup` scripts.
func ChezmoiApply() error {
	fmt.Println("chezmoi executed me")
	return nil
}
