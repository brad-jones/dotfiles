package updater

import (
	"syscall"
	"unsafe"

	"github.com/brad-jones/goerr/v2"
)

func hideFile(path string) (err error) {
	defer goerr.Handle(func(e error) { err = e })

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setFileAttributes := kernel32.NewProc("SetFileAttributesW")

	p, err := syscall.UTF16PtrFromString(path)
	goerr.Check(err, "failed to create pointer")

	r1, _, err := setFileAttributes.Call(uintptr(unsafe.Pointer(p)), 2)

	if r1 == 0 {
		goerr.Check(err)
	}

	return
}
