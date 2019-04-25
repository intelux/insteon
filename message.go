package insteon

import (
	"fmt"
	"strings"
)

// MessageFlags represents the message flags.
type MessageFlags byte

const (
	// MessageFlagExtended indicates extended messages.
	MessageFlagExtended MessageFlags = 0x10
	// MessageFlagAck indicates an acquitement message.
	MessageFlagAck MessageFlags = 0x20
	// MessageFlagAllLink indicates an all-link message.
	MessageFlagAllLink MessageFlags = 0x40
	// MessageFlagBroadcast indicates a broadcast message.
	MessageFlagBroadcast MessageFlags = 0x80
)

func (f MessageFlags) String() string {
	var result []string

	if f&MessageFlagExtended > 0 {
		result = append(result, "extended")
	}
	if f&MessageFlagAck > 0 {
		result = append(result, "ack")
	}
	if f&MessageFlagAllLink > 0 {
		result = append(result, "all-link")
	}
	if f&MessageFlagBroadcast > 0 {
		result = append(result, "broadcast")
	}

	return strings.Join(result, ",")
}

// Message is sent through the PLM to communicate with other devices.
type Message struct {
	Source       ID
	Target       ID
	HopsLeft     int
	MaxHops      int
	Flags        MessageFlags
	CommandBytes [2]byte
	UserData     [14]byte
}

func newMessage(target ID, commandBytes [2]byte) *Message {
	return &Message{
		Target:       target,
		HopsLeft:     2,
		MaxHops:      2,
		Flags:        0x00,
		CommandBytes: commandBytes,
	}
}

func newExtendedMessage(target ID, commandBytes [2]byte, userData [14]byte) *Message {
	return &Message{
		Target:       target,
		HopsLeft:     2,
		MaxHops:      2,
		Flags:        MessageFlagExtended,
		CommandBytes: commandBytes,
		UserData:     userData,
	}
}

// IsAck returns whether the message is an ack.
func (m Message) IsAck() bool {
	return m.Flags&MessageFlagAck != 0
}

// IsExtended returns whether the message is an extended message.
func (m Message) IsExtended() bool {
	return m.Flags&MessageFlagExtended != 0
}

// MarshalBinary -
func (m Message) MarshalBinary() ([]byte, error) {
	data := make([]byte, 9)

	copy(data[0:3], m.Source[:])
	copy(data[3:6], m.Target[:])

	flagsByte := byte(
		(m.MaxHops & 0x03) | (m.HopsLeft&0x03)<<2 | int(m.Flags),
	)

	data[6] = flagsByte
	copy(data[7:9], m.CommandBytes[:])

	if m.IsExtended() {
		data = append(data, m.UserData[:]...)
		data[len(data)-1] = checksum(m.CommandBytes, m.UserData[:])
	}

	return data, nil
}

// UnmarshalBinary -
func (m *Message) UnmarshalBinary(b []byte) error {
	if len(b) != 9 && len(b) != 23 {
		return fmt.Errorf("expected 9 or 23 bytes, got %d", len(b))
	}

	copy(m.Source[:], b[0:3])
	copy(m.Target[:], b[3:6])

	flagsByte := b[6]
	m.MaxHops = int(flagsByte) & 0x03
	m.HopsLeft = (int(flagsByte) & 0x0c) >> 2
	m.Flags = MessageFlags(flagsByte & 0xf0)

	copy(m.CommandBytes[:], b[7:9])

	if m.IsExtended() {
		if len(b) != 23 {
			return fmt.Errorf("message has the extended flag but not the expected size")
		}

		copy(m.UserData[:], b[9:23])
	} else {
		if len(b) != 9 {
			return fmt.Errorf("message does not have the extended flag but has the size")
		}
	}

	return nil
}

func checksum(commandBytes [2]byte, b []byte) byte {
	checksum := commandBytes[0] + commandBytes[1]

	for _, x := range b {
		checksum += x
	}

	return ((0xff ^ checksum) + 1) & 0xff
}
