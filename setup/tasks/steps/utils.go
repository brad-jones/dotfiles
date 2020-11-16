package steps

import (
	"os"
	"os/exec"

	"github.com/brad-jones/goexec/v2"
)

func isRoot() bool {
	return goexec.MustRunBuffered("id", "-u").StdOut == "0"
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
