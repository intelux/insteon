package plm

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"sync"

	"github.com/intelux/insteon/serial"
)

// PowerLineModem represents an Insteon PowerLine Modem device, which can be
// connected locally or via a TCP socket.
type PowerLineModem struct {
	Device io.ReadWriteCloser
	once   sync.Once
	stop   chan struct{}
	read   chan []byte
	write  chan []byte
}

// New create a new PowerLineModem device.
func New(device string) (*PowerLineModem, error) {
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

	plm := &PowerLineModem{
		Device: dev,
		stop:   make(chan struct{}),
		read:   make(chan []byte, 10),
		write:  make(chan []byte, 10),
	}

	go plm.readLoop()
	go plm.writeLoop()

	return plm, nil
}

// Close the PowerLine Modem.
func (m *PowerLineModem) Close() {
	m.once.Do(func() {
		close(m.stop)
		close(m.write)
	})

	m.Device.Close()
}

func (m *PowerLineModem) readLoop() {
	for {
		select {
		case <-m.stop:
			close(m.read)
			return
		default:
			msg := make([]byte, 16)
			n, err := m.Device.Read(msg)

			if err != nil {
				return
			}

			m.read <- msg[:n]
		}
	}
}

func (m *PowerLineModem) writeLoop() {
	for {
		select {
		case <-m.stop:
			return
		case msg := <-m.write:
			for len(msg) > 0 {
				n, err := m.Device.Write(msg)

				if err != nil {
					return
				}

				msg = msg[:n-1]
			}
		}
	}
}

// GetInfo gets information about the PowerLine Modem.
func (m *PowerLineModem) GetInfo() (Info, error) {
	m.write <- []byte{byte(GetIMInfo)}

	data := <-m.read

	fmt.Println(data)

	return Info{}, nil
}
