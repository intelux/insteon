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
	p.Payload = make([]byte, len(b)-3)
	copy(p.Payload, b[2:len(b)-1])
	p.Ack = b[len(b)-1]

	return nil
}

// MarshalBinary -
func (p packet) MarshalBinary() ([]byte, error) {
	b := make([]byte, len(p.Payload)+2)
	b[0] = messageStart
	b[1] = byte(p.CommandCode)
	copy(b[2:len(p.Payload)+2], p.Payload)

	if p.Ack != 0 {
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
