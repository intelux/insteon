package serial

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"
)

type termios struct {
	iflag  uint32
	oflag  uint32
	cflag  uint32
	lflag  uint32
	cc     [19]byte
	_      [3]uint32
	ispeed uint32
	ospeed uint32
}

const ioctlTCSETS = 0x00005402

var plmTermios = termios{
	// 8 bits, no parity, 1 stop bit. Also ignore modem status line.
	cflag:  syscall.CLOCAL | syscall.CREAD | 0x00001000 | syscall.CS8,
	cc:     [19]byte{},
	ispeed: 19200,
	ospeed: 19200,
}

func init() {
	plmTermios.cc[syscall.VTIME] = 0x01
	plmTermios.cc[syscall.VMIN] = 0x01
}

func open(device string) (io.ReadWriteCloser, error) {
	f, err := os.OpenFile(device, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0600)

	if err != nil {
		return nil, fmt.Errorf("failed to open serial port: %s", err)
	}

	err = syscall.SetNonblock(int(f.Fd()), false)

	if err != nil {
		return nil, err
	}

	r, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		uintptr(ioctlTCSETS),
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
