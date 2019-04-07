package plm

import (
	"fmt"
	"io"
	"sync"
)

// Connecter represents a type that can connect to a writer.
type Connecter interface {
	Connect(io.Writer) error
}

// ConnectedPipe is a pipe that can be dynamically connected.
type ConnectedPipe struct {
	mutex  sync.Mutex
	reader *io.PipeReader
	writer *io.PipeWriter
}

// Write to the connected pipe.
//
// If the pipe isn't connected at the moment of the write, the write is faked
// and will not block.
func (p *ConnectedPipe) Write(b []byte) (int, error) {
	p.mutex.Lock()
	w := p.writer
	p.mutex.Unlock()

	if w != nil {
		n, err := w.Write(b)

		if err != io.ErrClosedPipe {
			return n, err
		}
	}

	// The pipe is not connected, pretend to write suceeded.
	return len(b), nil
}

// Connect to a writer until either the pipe or the writer is closed.
//
// Connect returns the error that caused the copy to fail.
func (p *ConnectedPipe) Connect(w io.Writer) error {
	p.mutex.Lock()
	p.reader, p.writer = io.Pipe()
	r := p.reader
	p.mutex.Unlock()

	fmt.Println("copying...")
	_, err := io.Copy(w, r)
	fmt.Println("copying done", err)

	p.Close()

	return err
}

// Close the connected pipe.
//
// Any pending write or connect is interrupted.
//
// Close may be called as many times as needed. Closing an unconnected pipe has
// no effect.
func (p *ConnectedPipe) Close() (err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.reader != nil {
		p.reader.Close()
		p.reader = nil
	}

	if p.writer != nil {
		p.writer.Close()
		p.writer = nil
	}

	return err
}
