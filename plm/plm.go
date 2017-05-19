package plm

import (
	"fmt"
	"io"
	"net"
	"net/url"

	"github.com/intelux/insteon/serial"
)

// PowerLineModem represents an Insteon PowerLine Modem device, which can be
// connected locally or via a TCP socket.
type PowerLineModem struct {
	Device io.ReadWriteCloser
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

	return &PowerLineModem{
		Device: dev,
	}, nil
}

// Close the PowerLine Modem.
func (m *PowerLineModem) Close() {
	m.Device.Close()
}

// GetInfo gets information about the PowerLine Modem.
func (m *PowerLineModem) GetInfo() (Info, error) {
	return Info{}, nil
}
