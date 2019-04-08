package insteon

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"sync"

	"github.com/jacobsa/go-serial/serial"
)

// PowerLineModem represnts a powerline modem.
type PowerLineModem struct {
	// Device is the underlying device to use to send and receive PLM commands.
	//
	// Can be a local serial port or a remote one (TCP).
	Device io.ReadWriteCloser

	once     sync.Once
	routines chan func()
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

// GetIMInfo gets information about the PowerLine Modem.
func (m *PowerLineModem) GetIMInfo(ctx context.Context) (imInfo *IMInfo, err error) {
	m.init()

	if err = m.execute(ctx, func() error {
		if err := m.writeMessage(cmdGetIMInfo, nil); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return
}

func (m *PowerLineModem) init() {
	m.once.Do(func() {
		m.routines = make(chan func(), 10)

		go func() {
			for routine := range m.routines {
				routine()
			}
		}()
	})
}

func (m *PowerLineModem) execute(ctx context.Context, fn func() error) error {
	ch := make(chan error, 1)

	// Wait until we can push the routine.
	select {
	case m.routines <- func() {
		ch <- fn()
	}:
		select {
		case err := <-ch:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	case <-ctx.Done():
		return ctx.Err()
	}
}

const (
	// messageStart is the marker at the beginning of commands.
	messageStart byte = 0x02
	// messageAck is returned as an acknowledgment.
	messageAck byte = 0x06
	// messageNak is returned as an non-acknowledgment.
	messageNak byte = 0x15
)

func (m *PowerLineModem) writeMessage(commandCode CommandCode, msg interface{}) error {
	// Make sure we send it all at once. Not that it is mandatory, but it's
	// easier to debug.
	bw := bufio.NewWriter(m.Device)
	defer bw.Flush()

	_, err := bw.Write([]byte{
		messageStart,
		byte(commandCode),
	})

	if err != nil {
		return err
	}

	return NewMessageEncoder(bw).Encode(msg)
}
