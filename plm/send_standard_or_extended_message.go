package plm

import (
	"io"
)

// SendStandardOrExtendedMessageRequest is sent when information about is PLM is requested.
type SendStandardOrExtendedMessageRequest struct {
	Target       Identity
	HopsLeft     int
	MaxHops      int
	Flags        MessageFlags
	CommandBytes CommandBytes
	UserData     [14]byte
}

func (SendStandardOrExtendedMessageRequest) commandCode() CommandCode {
	return SendStandardOrExtendedMessage
}

func (r SendStandardOrExtendedMessageRequest) checksum() byte {
	return checksum(r.CommandBytes, r.UserData)
}

func (r SendStandardOrExtendedMessageRequest) marshal(w io.Writer) error {
	flagsByte := byte(
		(r.MaxHops & 0x03) | (r.HopsLeft&0x03)<<2 | int(r.Flags),
	)

	var data []byte
	data = append(data, r.Target[:]...)
	data = append(data, flagsByte)
	data = append(data, r.CommandBytes[:]...)

	if r.Flags&MessageFlagExtended != 0 {
		data = append(data, r.UserData[:]...)
		data[len(data)-1] = r.checksum()
	}

	_, err := w.Write(data)

	return err
}

// SendStandardOrExtendedMessageResponse is returned when information about is PLM is requested.
type SendStandardOrExtendedMessageResponse struct {
	Sender       Identity
	HopsLeft     int
	MaxHops      int
	Flags        MessageFlags
	CommandBytes CommandBytes
	UserData     [14]byte
}

func (*SendStandardOrExtendedMessageResponse) commandCode() CommandCode {
	return SendStandardOrExtendedMessage
}

func (res *SendStandardOrExtendedMessageResponse) unmarshal(r io.Reader) error {
	_, err := io.ReadFull(r, res.Sender[:])

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

	if res.Flags&MessageFlagExtended != 0 {
		_, err = io.ReadFull(r, res.UserData[:])

		if err != nil {
			return err
		}
	}

	_, err = io.ReadFull(r, data)

	if err != nil {
		return err
	}

	if data[0] != MessageAck {
		return ErrCommandFailure
	}

	return nil
}
