package insteon

import (
	"bufio"
	"io"
)

// PacketWriter implements a packet reader on top of a regular reader.
type packetWriter struct {
	writer *bufio.Writer
}

func newPacketWriter(r io.Writer) packetWriter {
	return packetWriter{bufio.NewWriter(r)}
}

func (w packetWriter) WritePacket(p *packet) error {
	b, err := p.MarshalBinary()

	if err != nil {
		return err
	}

	if _, err = w.writer.Write(b); err != nil {
		return err
	}

	return w.writer.Flush()
}
