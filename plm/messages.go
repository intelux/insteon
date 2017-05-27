package plm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

// CommandCoder is the interface for objects that have a command code.
type CommandCoder interface {
	commandCode() CommandCode
}

// Request is the interface for all requests.
type Request interface {
	CommandCoder
	write(io.Writer) error
}

// Response is the interface for all responses.
type Response interface {
	CommandCoder
	unmarshal(io.Reader) error
}

var (
	// ErrCommandFailure indicates that a previous command failed.
	ErrCommandFailure = errors.New("command failed")
)

type errUnknownCommand struct {
	CommandCode
}

// Error returns the error string.
func (e errUnknownCommand) Error() string {
	return fmt.Sprintf("unknown command code (%x)", e.CommandCode)
}

// MarshalRequest serializes a request to a writer.
func MarshalRequest(w io.Writer, request Request) error {
	// Make sure we send it all at once. Not that it is mandatory, but it's
	// easier to debug.
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	_, err := bw.Write([]byte{
		byte(MessageStart),
		byte(request.commandCode()),
	})

	if err != nil {
		return err
	}

	return request.write(bw)
}

// UnmarshalResponse parses a response from a reader.
func UnmarshalResponse(r io.Reader, response Response) error {
	for {
		mark, err := skipToMessage(r)

		if err != nil {
			return err
		}

		if mark == MessageNak {
			return ErrCommandFailure
		}

		buffer := make([]byte, 1)

		_, err = r.Read(buffer)

		if err != nil {
			return err
		}

		commandCode := CommandCode(buffer[0])

		if commandCode != response.commandCode() {
			buffer := make([]byte, responsesSizes[commandCode])
			_, err = io.ReadFull(r, buffer)

			if err != nil {
				return err
			}
		} else {
			err = response.unmarshal(r)

			if err != nil {
				return err
			}

			return nil
		}
	}
}

// skipToMessage skips bytes on the specified io.Reader until a message start
// is met.
//
// The message start or message nak is returned.
func skipToMessage(r io.Reader) (byte, error) {
	buffer := make([]byte, 1)

	for buffer[0] != MessageStart && buffer[0] != MessageNak {
		_, err := r.Read(buffer)

		if err != nil {
			return 0, err
		}
	}

	return buffer[0], nil
}

var responsesSizes = map[CommandCode]int{
	GetIMInfo:                     7,
	SendStandardOrExtendedMessage: 7,
	StandardMessageReceived:       9,
	ExtendedMessageReceived:       23,
}
