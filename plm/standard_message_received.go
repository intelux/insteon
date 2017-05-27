package plm

import (
	"io"
)

// StandardMessageReceivedResponse is returned when the PLM sends an unsollicited
// message to the host.
type StandardMessageReceivedResponse struct {
	Sender       Identity
	Target       Identity
	HopsLeft     int
	MaxHops      int
	Flags        MessageFlags
	CommandBytes CommandBytes
}

func (*StandardMessageReceivedResponse) commandCode() CommandCode {
	return StandardMessageReceived
}

func (res *StandardMessageReceivedResponse) unmarshal(r io.Reader) error {
	_, err := io.ReadFull(r, res.Sender[:])

	if err != nil {
		return err
	}

	_, err = io.ReadFull(r, res.Target[:])

	if err != nil {
		return err
	}

	data := make([]byte, 1)
	_, err = io.ReadFull(r, data)

	if err != nil {
		return err
	}

	flagsByte := data[0]
	res.MaxHops = int(flagsByte) & 0x03
	res.HopsLeft = (int(flagsByte) & 0x0c) >> 2
	res.Flags = MessageFlags(flagsByte & 0xf0)

	_, err = io.ReadFull(r, res.CommandBytes[:])

	if err != nil {
		return err
	}

	return nil
}
