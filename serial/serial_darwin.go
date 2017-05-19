package serial

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"
)

type termios struct {
	iflag  uint64
	oflag  uint64
	cflag  uint64
	lflag  uint64
	cc     [20]byte
	ispeed uint64
	ospeed uint64
}

const ioctlTIOCSETA = 0x80487414

var plmTermios = termios{
	// 8 bits, no parity, 1 stop bit. Also ignore modem status line.
	cflag: 0x00008B00,
	cc: [20]byte{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		// Read no less than 1 character each time and wait at most 1 tenth of
		// second after each character.
		1, 1, 0, 0,
	},
	ispeed: 19200,
	ospeed: 19200,
}

func open(device string) (io.ReadWriteCloser, error) {
	f, err := os.OpenFile(device, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0600)

	if err != nil {
		return nil, fmt.Errorf("failed to open serial port: %s", err)
	}

	r, _, errno := syscall.Syscall(
		syscall.SYS_FCNTL,
		uintptr(f.Fd()),
		uintptr(syscall.F_SETFL),
		uintptr(0),
	)

	if errno != 0 {
		return nil, os.NewSyscallError("SYS_FCNTL", errno)
	}

	if r != 0 {
		return nil, fmt.Errorf("SYS_FCNTL call returned %d", r)
	}

	r, _, errno = syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		uintptr(ioctlTIOCSETA),
		uintptr(unsafe.Pointer(&plmTermios)),
	)

	if errno != 0 {
		return nil, os.NewSyscallError("SYS_IOCTL", errno)
	}

	// Just in case, check the return value as well.
	if r != 0 {
		return nil, fmt.Errorf("SYS_IOCTL call returned %d", r)
	}

	return f, nil
}
