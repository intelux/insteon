package insteon

import (
	"context"
	"encoding"
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

	once            sync.Once
	ctx             context.Context
	cancel          func()
	routines        chan func()
	incomingPackets chan *packet
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

	err = m.execute(ctx, func(ctx context.Context) error {
		imInfo = &IMInfo{}

		return m.transfer(ctx, &packet{CommandCode: cmdGetIMInfo}, imInfo)
	})

	return
}

func (m *PowerLineModem) init() {
	m.once.Do(func() {
		m.ctx, m.cancel = context.WithCancel(context.Background())
		m.routines = make(chan func())
		m.incomingPackets = make(chan *packet)

		go m.readLoop(m.ctx)

		go func() {
			for routine := range m.routines {
				routine()
			}
		}()
	})
}

func (m *PowerLineModem) execute(ctx context.Context, fn func(context.Context) error) error {
	ch := make(chan error, 1)

	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	go func() {
		<-m.ctx.Done()
		cancel()
	}()

	// Wait until we can push the routine.
	select {
	case m.routines <- func() {
		ch <- fn(ctx)
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

func (m *PowerLineModem) readLoop(ctx context.Context) {
	r := newPacketReader(m.Device)

	for {
		p, err := r.ReadPacket()

		if err != nil {
			return
		}

		select {
		case m.incomingPackets <- p:
		case <-ctx.Done():
			return
		}
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

func (m *PowerLineModem) readPacket(ctx context.Context, commandCode CommandCode) (*packet, error) {
	for {
		select {
		case packet := <-m.incomingPackets:
			if packet.CommandCode == commandCode {
				return packet, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (m *PowerLineModem) writePacket(p *packet) error {
	w := newPacketWriter(m.Device)

	return w.WritePacket(p)
}

func (m *PowerLineModem) transfer(ctx context.Context, p *packet, result encoding.BinaryUnmarshaler) error {
	if err := m.writePacket(p); err != nil {
		return err
	}

	p, err := m.readPacket(ctx, p.CommandCode)

	if err != nil {
		return err
	}

	return result.UnmarshalBinary(p.Payload)
}
