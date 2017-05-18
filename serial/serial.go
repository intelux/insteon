package serial

import "io"

// Open a new serial port from the specified device.
//
// The device is assumed to be a PLM and will be configured accordingly.
func Open(device string) (io.ReadWriteCloser, error) {
	return open(device)
}
