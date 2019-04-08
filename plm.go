package insteon

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"os"

	"github.com/jacobsa/go-serial/serial"
)

// PowerLineModem represnts a powerline modem.
type PowerLineModem struct {
	// Device is the underlying device to use to send and receive PLM commands.
	//
	// Can be a local serial port or a remote one (TCP).
	Device io.ReadWriteCloser
}

// DefaultPowerLineModem is the default PowerLine Modem instance.
var DefaultPowerLineModem = func() *PowerLineModem {
	plm, err := NewPowerLineModem(PowerLineModemDevice)

	if err != nil {
		panic(fmt.Errorf("instanciating default PowerLine Modem: %s", err))
	}

	return plm
}()

// NewLocalPowerLineModem instantiates a new local PowerLine Modem.
func NewLocalPowerLineModem(serialPort string) (*PowerLineModem, error) {
	options := serial.OpenOptions{
		PortName:        serialPort,
		BaudRate:        19200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	// Open the port.
	var device io.ReadWriteCloser
	var err error

	if device, err = serial.Open(options); err != nil {
		return nil, fmt.Errorf("opening local serial port: %s", err)
	}

	if PowerLineModemDebug {
		device = debugReadWriteCloser{
			ReadWriteCloser: device,
			DebugWriter:     os.Stderr,
		}
	}

	return &PowerLineModem{
		Device: device,
	}, nil
}

// NewRemotePowerLineModem instantiates a new remote PowerLine Modem.
func NewRemotePowerLineModem(host string) (*PowerLineModem, error) {
	var device io.ReadWriteCloser
	var err error

	if device, err = net.Dial("tcp", host); err != nil {
		return nil, fmt.Errorf("opening remote serial port: %s", err)
	}

	if PowerLineModemDebug {
		device = debugReadWriteCloser{
			ReadWriteCloser: device,
			DebugWriter:     os.Stderr,
		}
	}

	return &PowerLineModem{
		Device: device,
	}, nil
}

// NewPowerLineModem instantiates a new PowerLine Modem.
func NewPowerLineModem(device string) (*PowerLineModem, error) {
	url, err := url.Parse(device)

	if err != nil {
		return nil, fmt.Errorf("parsing device: %s", err)
	}

	switch url.Scheme {
	case "tcp":
		return NewRemotePowerLineModem(url.Host)
	default:
		return NewLocalPowerLineModem(url.String())
	}
}
