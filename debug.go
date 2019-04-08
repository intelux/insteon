package insteon

import (
	"encoding/hex"
	"fmt"
	"io"
)

type debugReadWriteCloser struct {
	io.ReadWriteCloser
	DebugWriter io.Writer
}

func (d debugReadWriteCloser) Read(buf []byte) (n int, err error) {
	defer func() {
		if err == nil {
			fmt.Fprintf(d.DebugWriter, "< %s\n", hex.EncodeToString(buf[:n]))
		} else {
			fmt.Fprintf(d.DebugWriter, "< %s\n", err)
		}
	}()

	return d.ReadWriteCloser.Read(buf)
}

// Write to the device.
func (d debugReadWriteCloser) Write(buf []byte) (n int, err error) {
	defer func() {
		if err == nil {
			fmt.Fprintf(d.DebugWriter, "> %s\n", hex.EncodeToString(buf[:n]))
		} else {
			fmt.Fprintf(d.DebugWriter, "> %s\n", err)
		}
	}()

	return d.ReadWriteCloser.Write(buf)
}
