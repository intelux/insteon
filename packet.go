package insteon

import "fmt"

type packet struct {
	CommandCode CommandCode
	Payload     []byte
	Ack         byte
}

// UnmarshalBinary -
func (p *packet) UnmarshalBinary(b []byte) error {
	if b[0] != messageStart {
		return fmt.Errorf("expected %02x but got %02x", messageStart, b[0])
	}

	p.CommandCode = CommandCode(b[1])

	// Always treat outgoing messages as if they started with the sender
	// identity so that they can be parsed as messages.
	if p.CommandCode == cmdSendStandardOrExtendedMessage {
		extB := make([]byte, len(b)+3)
		copy(extB[3:], b)
		b = extB
	}

	p.Payload = make([]byte, len(b)-2)
	copy(p.Payload, b[2:len(b)])

	if isOutgoingCommandCode(p.CommandCode) {
		p.Ack = p.Payload[len(p.Payload)-1]
		p.Payload = p.Payload[:len(p.Payload)-1]
	}

	return nil
}

// MarshalBinary -
func (p packet) MarshalBinary() ([]byte, error) {
	payload := p.Payload

	// Strip the sender identity for outgoing messages as it must not really be
	// specified.
	if p.CommandCode == cmdSendStandardOrExtendedMessage {
		payload = payload[3:]
	}

	b := make([]byte, len(payload)+2)
	b[0] = messageStart
	b[1] = byte(p.CommandCode)
	copy(b[2:len(payload)+2], payload)

	if isOutgoingCommandCode(p.CommandCode) && p.Ack != 0 {
		b = append(b, p.Ack)
	}

	return b, nil
}

func (p packet) IsAck() bool {
	return p.Ack == messageAck
}

func (p packet) IsNak() bool {
	return p.Ack == messageNak
}
