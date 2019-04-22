package insteon

import (
	"context"
	"encoding"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

// SerialPowerLineModem represnts a powerline modem.
type SerialPowerLineModem struct {
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

// NewLocalPowerLineModem instantiates a new local PowerLine Modem.
func NewLocalPowerLineModem(serialPort string) (*SerialPowerLineModem, error) {
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

	return &SerialPowerLineModem{
		Device: device,
	}, nil
}

// NewRemotePowerLineModem instantiates a new remote PowerLine Modem.
func NewRemotePowerLineModem(host string) (*SerialPowerLineModem, error) {
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

	return &SerialPowerLineModem{
		Device: device,
	}, nil
}

// GetIMInfo gets information about the PowerLine Modem.
func (m *SerialPowerLineModem) GetIMInfo(ctx context.Context) (imInfo *IMInfo, err error) {
	m.init()

	err = m.execute(ctx, func(ctx context.Context) error {
		imInfo = &IMInfo{}

		return m.roundtrip(ctx, &packet{CommandCode: cmdGetIMInfo}, imInfo)
	})

	return
}

// GetAllLinkDB gets the on level of a device.
func (m *SerialPowerLineModem) GetAllLinkDB(ctx context.Context) (records AllLinkRecordSlice, err error) {
	m.init()

	err = m.execute(ctx, func(ctx context.Context) error {
		p, err := m.rawRoundtrip(ctx, &packet{CommandCode: cmdGetFirstAllLinkRecord})

		if err != nil {
			return err
		}

		// A NAK at this point indicates that the DB is empty.
		if p.IsNak() {
			return nil
		}

		record := &AllLinkRecord{}

		if _, err := m.readPacketTo(ctx, cmdAllLinkRecordMessage, record); err != nil {
			return err
		}

		records = append(records, *record)

		for {
			p, err := m.rawRoundtrip(ctx, &packet{CommandCode: cmdGetNextAllLinkRecord})

			if err != nil {
				break
			}

			// A NAK at this point indicates that the listing is over.
			if p.IsNak() {
				break
			}

			if _, err := m.readPacketTo(ctx, cmdAllLinkRecordMessage, record); err != nil {
				return err
			}

			records = append(records, *record)
		}

		return nil
	})

	sort.Stable(records)

	return
}

// GetDeviceState gets the on level of a device.
func (m *SerialPowerLineModem) GetDeviceState(ctx context.Context, identity ID) (state *LightState, err error) {
	m.init()

	err = m.execute(ctx, func(ctx context.Context) error {
		msg := newMessage(identity, commandBytesStatusRequest)
		_, err := m.messageRoundtrip(ctx, msg)

		if err != nil {
			return err
		}

		rmsg, err := m.readStandardMessage(ctx)

		if err != nil {
			return err
		}

		level := byteToOnLevel(rmsg.CommandBytes[1])

		state = &LightState{
			OnOff: level > 0,
			Level: level,
		}

		return nil
	})

	return
}

// SetDeviceState sets the state of a lighting device.
func (m *SerialPowerLineModem) SetDeviceState(ctx context.Context, identity ID, state LightState) (err error) {
	m.init()

	err = m.execute(ctx, func(ctx context.Context) error {
		msg := newMessage(identity, state.asCommandBytes())
		_, err := m.messageRoundtrip(ctx, msg)

		return err
	})

	return
}

// GetDeviceInfo returns the information about a device.
func (m *SerialPowerLineModem) GetDeviceInfo(ctx context.Context, identity ID) (deviceInfo *DeviceInfo, err error) {
	m.init()

	err = m.execute(ctx, func(ctx context.Context) error {
		msg := newExtendedMessage(identity, commandBytesGetDeviceInfo, [14]byte{})
		_, err := m.messageRoundtrip(ctx, msg)

		if err != nil {
			return err
		}

		rmsg, err := m.readExtendedMessage(ctx)

		if err != nil {
			return err
		}

		deviceInfo = &DeviceInfo{}

		return deviceInfo.UnmarshalBinary(rmsg.UserData[:])
	})

	return
}

// SetDeviceInfo sets the information on device.
func (m *SerialPowerLineModem) SetDeviceInfo(ctx context.Context, identity ID, deviceInfo DeviceInfo) (err error) {
	if deviceInfo.X10Address != nil {
		if err = m.SetDeviceX10Address(ctx, identity, *deviceInfo.X10Address); err != nil {
			return err
		}
	}

	if deviceInfo.RampRate != nil {
		if err = m.SetDeviceRampRate(ctx, identity, *deviceInfo.RampRate); err != nil {
			return err
		}
	}

	if deviceInfo.OnLevel != nil {
		if err = m.SetDeviceOnLevel(ctx, identity, *deviceInfo.OnLevel); err != nil {
			return err
		}
	}

	if deviceInfo.LEDBrightness != nil {
		if err = m.SetDeviceLEDBrightness(ctx, identity, *deviceInfo.LEDBrightness); err != nil {
			return err
		}
	}

	return
}

// SetDeviceX10Address sets a device X10 address.
func (m *SerialPowerLineModem) SetDeviceX10Address(ctx context.Context, identity ID, x10Address [2]byte) (err error) {
	m.init()

	err = m.execute(ctx, func(ctx context.Context) error {
		userData := [14]byte{}
		userData[1] = 0x04
		userData[2] = x10Address[0]
		userData[3] = x10Address[1]

		msg := newExtendedMessage(identity, commandBytesSetDeviceInfo, userData)
		_, err := m.messageRoundtrip(ctx, msg)

		return err
	})

	return
}

// SetDeviceRampRate sets a device ramp rate.
func (m *SerialPowerLineModem) SetDeviceRampRate(ctx context.Context, identity ID, rampRate time.Duration) (err error) {
	m.init()

	err = m.execute(ctx, func(ctx context.Context) error {
		userData := [14]byte{}
		userData[1] = 0x05
		userData[2] = rampRateToByte(rampRate)

		msg := newExtendedMessage(identity, commandBytesSetDeviceInfo, userData)
		_, err := m.messageRoundtrip(ctx, msg)

		return err
	})

	return
}

// SetDeviceOnLevel sets a device on level.
func (m *SerialPowerLineModem) SetDeviceOnLevel(ctx context.Context, identity ID, level float64) (err error) {
	m.init()

	err = m.execute(ctx, func(ctx context.Context) error {
		userData := [14]byte{}
		userData[1] = 0x06
		userData[2] = onLevelToByte(level)

		msg := newExtendedMessage(identity, commandBytesSetDeviceInfo, userData)
		_, err := m.messageRoundtrip(ctx, msg)

		return err
	})

	return
}

// SetDeviceLEDBrightness sets a device LED brightness.
func (m *SerialPowerLineModem) SetDeviceLEDBrightness(ctx context.Context, identity ID, level float64) (err error) {
	m.init()

	err = m.execute(ctx, func(ctx context.Context) error {
		userData := [14]byte{}
		userData[1] = 0x07
		userData[2] = ledBrightnessToByte(level)

		msg := newExtendedMessage(identity, commandBytesSetDeviceInfo, userData)
		_, err := m.messageRoundtrip(ctx, msg)

		return err
	})

	return
}

// Beep causes a device to beep.
func (m *SerialPowerLineModem) Beep(ctx context.Context, identity ID) (err error) {
	m.init()

	err = m.execute(ctx, func(ctx context.Context) error {
		msg := newMessage(identity, commandBytesBeep)
		_, err := m.messageRoundtrip(ctx, msg)

		return err
	})

	return
}

func (m *SerialPowerLineModem) init() {
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

func (m *SerialPowerLineModem) execute(ctx context.Context, fn func(context.Context) error) error {
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

func (m *SerialPowerLineModem) readLoop(ctx context.Context) {
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

func (m *SerialPowerLineModem) readPacket(ctx context.Context, commandCode CommandCode) (*packet, error) {
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

func (m *SerialPowerLineModem) readPacketTo(ctx context.Context, commandCode CommandCode, result encoding.BinaryUnmarshaler) (*packet, error) {
	p, err := m.readPacket(ctx, commandCode)

	if err != nil {
		return nil, err
	}

	if result != nil {
		return p, result.UnmarshalBinary(p.Payload)
	}

	return p, nil
}

func (m *SerialPowerLineModem) writePacket(p *packet) error {
	w := newPacketWriter(m.Device)

	return w.WritePacket(p)
}

func (m *SerialPowerLineModem) messageRoundtrip(ctx context.Context, msg *Message) (*Message, error) {
	payload, err := msg.MarshalBinary()

	if err != nil {
		return nil, fmt.Errorf("marshalling message: %s", err)
	}

	p := &packet{
		CommandCode: cmdSendStandardOrExtendedMessage,
		Payload:     payload,
	}

	result := &Message{}

	if err = m.roundtrip(ctx, p, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (m *SerialPowerLineModem) readMessage(ctx context.Context, commandCode CommandCode) (*Message, error) {
	result := &Message{}

	if _, err := m.readPacketTo(ctx, commandCode, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (m *SerialPowerLineModem) readStandardMessage(ctx context.Context) (*Message, error) {
	return m.readMessage(ctx, cmdStandardMessageReceived)
}

func (m *SerialPowerLineModem) readExtendedMessage(ctx context.Context) (*Message, error) {
	return m.readMessage(ctx, cmdExtendedMessageReceived)
}

func (m *SerialPowerLineModem) rawRoundtrip(ctx context.Context, p *packet) (*packet, error) {
	if err := m.writePacket(p); err != nil {
		return nil, err
	}

	return m.readPacket(ctx, p.CommandCode)
}

func (m *SerialPowerLineModem) roundtrip(ctx context.Context, p *packet, result encoding.BinaryUnmarshaler) (err error) {
	var rp *packet

	for {
		rp, err = m.rawRoundtrip(ctx, p)

		if err != nil {
			return err
		}

		if rp.IsAck() {
			break
		}

		select {
		case <-time.After(time.Millisecond * 150):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if result != nil {
		return result.UnmarshalBinary(rp.Payload)
	}

	return nil
}
