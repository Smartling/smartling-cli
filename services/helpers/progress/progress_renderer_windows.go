//go:build windows

package progress

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procGetConsoleScreenBufferInfo = kernel32.NewProc(
		"GetConsoleScreenBufferInfo",
	)

	procSetConsoleCursorPosition = kernel32.NewProc(
		"SetConsoleCursorPosition",
	)
)

type (
	consoleCoordinates struct {
		X int16
		Y int16
	}

	consoleScreenBufferInfo struct {
		size              consoleCoordinates
		cursorPosition    consoleCoordinates
		maximumWindowSize consoleCoordinates
		attributes        int16
		window            struct {
			left   int16
			top    int16
			right  int16
			bottom int16
		}
	}
)

type ProgressRenderer struct{}

func (renderer ProgressRenderer) Render(progress *Progress) error {
	var info consoleScreenBufferInfo

	_, _, code := syscall.SyscallN(
		procGetConsoleScreenBufferInfo.Addr(),
		2,
		uintptr(syscall.Stderr),
		uintptr(unsafe.Pointer(&info)),
		0,
	)
	if code != 0 {
		return error(code)
	}

	pos := info.cursorPosition
	pos.X = 0

	_, _, code = syscall.SyscallN(
		procSetConsoleCursorPosition.Addr(),
		2,
		uintptr(syscall.Stderr),
		uintptr(uint32(uint16(pos.Y))<<16|uint32(uint16(pos.X))),
		0,
	)
	if code != 0 {
		return error(code)
	}

	_, err := fmt.Fprint(os.Stderr, progress.String())

	return err
}
