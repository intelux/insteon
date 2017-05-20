package plm

import (
	"encoding/hex"
	"fmt"
	"io"
)

type debugStream struct {
	io.Writer
	Local  string
	Remote string
}

func (s debugStream) Log(direction string, buf []byte, err error) {
	if err != nil {
		fmt.Fprintf(s, "[%s %s %s] %s\n", s.Local, direction, s.Remote, err)
	} else {
		fmt.Fprintf(
			s, "[%s %s %s] %s\n", s.Local, direction, s.Remote,
			hex.EncodeToString(buf),
		)
	}
}

// DebugReader is a io.Reader wrapper that outputs information about its reads.
type DebugReader struct {
	io.Reader
	debugStream
}

// Read from the device.
func (d DebugReader) Read(buf []byte) (n int, err error) {
	defer func() {
		d.Log("<", buf[:n], err)
	}()

	return d.Reader.Read(buf)
}

// DebugWriter is a io.Writer wrapper that outputs information about its
// writes.
type DebugWriter struct {
	io.Writer
	debugStream
}

// Write to the device.
func (d DebugWriter) Write(buf []byte) (n int, err error) {
	defer func() {
		d.Log(">", buf[:n], err)
	}()

	return d.Writer.Write(buf)
}
