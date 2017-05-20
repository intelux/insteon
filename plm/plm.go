package plm

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/url"

	"github.com/intelux/insteon/serial"
)

type requestToken struct {
	*io.PipeReader
	io.Writer
	pipeWriter *io.PipeWriter
	ready      chan struct{}
}

// Close the token.
func (t *requestToken) Close() error {
	t.PipeReader.Close()
	t.pipeWriter.Close()

	return nil
}

// PowerLineModem represents an Insteon PowerLine Modem device, which can be
// connected locally or via a TCP socket.
type PowerLineModem struct {
	reader         io.Reader
	writer         io.Writer
	closer         io.Closer
	dispatchReader *io.PipeReader
	dispatchWriter *io.PipeWriter
	stop           chan struct{}
	tokens         chan *requestToken
}

// ParseDevice parses a device specifiction string, either as a local file (to
// a serial port likely) or as a tcp:// URL.
func ParseDevice(device string) (io.ReadWriteCloser, error) {
	var err error
	url, _ := url.Parse(device)
	var dev io.ReadWriteCloser

	switch url.Scheme {
	case "tcp":
		dev, err = net.Dial("tcp", url.Host)

		if err != nil {
			return nil, fmt.Errorf("failed to connect to TCP device: %s", err)
		}
	case "":
		dev, err = serial.Open(url.String())

		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported scheme for device `%s`", url.Scheme)
	}

	return dev, nil
}

// New create a new PowerLineModem device.
func New(device io.ReadWriteCloser) *PowerLineModem {
	r, w := io.Pipe()

	return &PowerLineModem{
		reader:         io.TeeReader(device, w),
		writer:         device,
		closer:         device,
		dispatchReader: r,
		dispatchWriter: w,
	}
}

// SetDebugStream enables debug output on the specified writer.
func (m *PowerLineModem) SetDebugStream(w io.Writer) {
	debugStream := debugStream{
		Writer: w,
		Local:  "host",
		Remote: "plm",
	}

	m.reader = DebugReader{
		Reader:      m.reader,
		debugStream: debugStream,
	}
	m.writer = DebugWriter{
		Writer:      m.writer,
		debugStream: debugStream,
	}
}

// Start the PowerLine Modem.
//
// Attempting to start an already running intance has undefined behavior.
func (m *PowerLineModem) Start() {
	m.stop = make(chan struct{})
	go readLoop(m.stop, m.reader)

	m.tokens = make(chan *requestToken)
	go dispatchLoop(m.tokens, m.dispatchReader)
}

// Stop the PowerLine Modem.
//
// Attempting to stop a non-running intance has undefined behavior.
func (m *PowerLineModem) Stop() {
	close(m.tokens)
	m.tokens = nil

	close(m.stop)
	m.stop = nil
}

// Close the PowerLine Modem.
func (m *PowerLineModem) Close() {
	m.Stop()
	m.dispatchReader.Close()
	m.dispatchWriter.Close()
	m.closer.Close()
}

func readLoop(stop <-chan struct{}, r io.Reader) {
	for {
		select {
		case <-stop:
			return
		default:
			msg := make([]byte, 16)
			_, err := r.Read(msg)

			if err != nil {
				return
			}

			// TODO: Handle the monitor events.
			//fmt.Println(msg[:n])
		}
	}
}

func dispatchLoop(tokens <-chan *requestToken, r io.Reader) {
	for token := range tokens {
		close(token.ready)
		_, err := io.Copy(token.pipeWriter, r)

		if err != nil {
			return
		}
	}
}

func (m *PowerLineModem) createToken() *requestToken {
	r, w := io.Pipe()

	token := &requestToken{
		PipeReader: r,
		Writer:     m.writer,
		pipeWriter: w,
		ready:      make(chan struct{}),
	}

	m.tokens <- token

	return token
}

// Acquire the PowerLine Modem for exclusive reading-writing.
//
// It is the responsibility of the caller to close the returned instance.
func (m *PowerLineModem) Acquire(ctx context.Context) (io.ReadWriteCloser, error) {
	token := m.createToken()

	select {
	case <-token.ready:
		return token, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetInfo gets information about the PowerLine Modem.
func (m *PowerLineModem) GetInfo(ctx context.Context) (Info, error) {
	token, err := m.Acquire(ctx)

	if err != nil {
		return Info{}, err
	}

	defer token.Close()

	_, err = token.Write([]byte{MessageStart, byte(GetIMInfo)})

	if err != nil {
		return Info{}, err
	}

	buf := make([]byte, 20)
	_, err = token.Read(buf)

	if err != nil {
		return Info{}, err
	}

	return Info{}, nil
}
