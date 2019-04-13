package insteon

import (
	"bufio"
	"io"
)

// PacketReader implements a packet reader on top of a regular reader.
type packetReader struct {
	reader *bufio.Reader
}

func newPacketReader(r io.Reader) packetReader {
	return packetReader{bufio.NewReader(r)}
}

func (r packetReader) Read() ([]byte, error) {
	var result []byte

	for {
		b, err := r.reader.ReadByte()

		if err != nil {
			return nil, err
		}

		// We skip all bytes until we reach a message start.
		if b != messageStart {
			continue
		}

		if b, err = r.reader.ReadByte(); err != nil {
			return nil, err
		}

		commandCode := CommandCode(b)

		size, ok := packetSizes[commandCode]

		if !ok {
			// We didn't find a known command-code: let's ignore this packet
			// and wait for the next packet.
			if err = r.reader.UnreadByte(); err != nil {
				return nil, err
			}
		}

		result = make([]byte, size+2)
		result[0] = messageStart
		result[1] = b

		if _, err = io.ReadAtLeast(r.reader, result[2:], size); err != nil {
			return nil, err
		}

		// We may have to handle extended packets.

		return result, nil
	}
}

func (r packetReader) ReadPacket() (*packet, error) {
	b, err := r.Read()

	if err != nil {
		return nil, err
	}

	result := &packet{}

	if err = result.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return result, nil
}

var packetSizes = map[CommandCode]int{
	cmdStandardMessageReceived:       9,
	cmdExtendedMessageReceived:       23,
	cmdX10Received:                   2,
	cmdAllLinkingCompleted:           8,
	cmdButtonEventReport:             1,
	cmdUserResetDetected:             0,
	cmdAllLinkCleanupFailureReport:   5,
	cmdAllLinkRecordMessage:          8,
	cmdAllLinkCleanupStatusReport:    1,
	cmdGetIMInfo:                     7,
	cmdGetFirstAllLinkRecord:         1,
	cmdGetNextAllLinkRecord:          1,
	cmdStartAllLinking:               3,
	cmdCancelAllLinking:              1,
	cmdSendStandardOrExtendedMessage: 7,
}
