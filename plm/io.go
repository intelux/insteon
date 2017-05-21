package plm

import (
	"io"
	"sync"
)

// ConnectReader represents a reader that can be connected.
type ConnectReader interface {
	io.Reader
	Connecter
}

// Connecter represents an object that can be connected.
type Connecter interface {
	Connect()
	Disconnect()
}

type connectWriter struct {
	io.Writer
	connected      bool
	connectedMutex sync.Mutex
}

func (w *connectWriter) Write(b []byte) (int, error) {
	w.connectedMutex.Lock()
	c := w.connected
	w.connectedMutex.Unlock()

	if c {
		return w.Writer.Write(b)
	}

	return len(b), nil
}

func (w *connectWriter) connect() {
	w.connectedMutex.Lock()
	w.connected = true
	w.connectedMutex.Unlock()
}

func (w *connectWriter) disconnect() {
	w.connectedMutex.Lock()
	w.connected = false
	w.connectedMutex.Unlock()
}

type connectReader struct {
	io.Reader
	writer *connectWriter
}

func (r connectReader) Connect() {
	r.writer.connect()
}

func (r connectReader) Disconnect() {
	r.writer.disconnect()
}

// ConnectPipe connects two ends of an io.Pipe.
func ConnectPipe(r io.Reader, w io.Writer) (ConnectReader, io.Writer) {
	cw := &connectWriter{
		Writer: w,
	}

	return connectReader{
		Reader: r,
		writer: cw,
	}, cw
}
